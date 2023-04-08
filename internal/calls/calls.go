package calls

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/eyeonicon/go-icon-sdk/transactions"
	"github.com/eyeonicon/go-icon-sdk/util"
	"github.com/icon-project/goloop/client"
	"sync"
	"time"
)

var BOOSTED_OMM = "cxeaff5a10cb72bf85965b8b4af3e708ab772b7921"
var DELEGATION = "cx841f29ec6ce98b527d49a275e87d427627f1afe5"

type User struct {
	address string
	votes   *big.Int
}

type VoteInfo struct {
	Address    string `json:"_address"`
	VotesInIcx string `json:"_votes_in_icx"`
	VotesInPer string `json:"_votes_in_per"`
}

// Get the addresses of all known stakers
func GetStakers(c *client.ClientV3) ([]string, error) {
	start := time.Now()

	var stakers []string
	amountOfUsers, err := getAmountOfOMMUsers(c)

	if err != nil {
		return nil, err
	}

	amountOfSkips := int(amountOfUsers/100 + 1) // - 12 // minus 12 is for testing

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i <= amountOfSkips; i++ {
		wg.Add(1)

		go func(_c *client.ClientV3, _i int) {
			defer wg.Done()
			users, err := getOMMUsers(c, _i)
			if err != nil {
				panic(err)
			}
			for _, user := range users {
				mu.Lock()
				stakers = append(stakers, user)
				mu.Unlock()
			}
		}(c, i)
	}

	wg.Wait()
	fmt.Printf("GetStakers took %.2f seconds\n", time.Since(start).Seconds())
	return stakers, nil
}

// returns amount of omm users
func getAmountOfOMMUsers(c *client.ClientV3) (int64, error) {
	callObj := transactions.CallBuilder(BOOSTED_OMM, "activeUsersCount", nil)
	res, err := c.Call(callObj)
	if err != nil {
		fmt.Println(err)
	}

	hexStr := fmt.Sprintf("%v", res)
	intVal, err := strconv.ParseInt(hexStr, 0, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	return intVal, nil
}

// returns all the omm users
func getOMMUsers(c *client.ClientV3, skip int) ([]string, error) {

	start := int64(0 + (skip * 100))
	end := int64(start + 100)

	params := map[string]interface{}{
		"start": "0x" + strconv.FormatInt(start, 16),
		"end":   "0x" + strconv.FormatInt(end, 16),
	}

	callObj := transactions.CallBuilder(BOOSTED_OMM, "getUsers", params)
	res, err := c.Call(callObj)

	if err != nil {
		panic(err)
		// return nil, err
	}

	strSlice, ok := res.([]string)
	if !ok {
		if interfSlice, ok := res.([]interface{}); ok {
			strSlice = make([]string, len(interfSlice))
			for i, v := range interfSlice {
				strSlice[i] = fmt.Sprint(v)
			}
		}
	}

	return strSlice, nil
}

// returns a list of all users's address and vote amount on validator
func GetValidatorVotes(c *client.ClientV3, validator string) []User {
	start := time.Now()

	var validatorVotes []User

	users, err := GetStakers(c)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, user := range users {
		wg.Add(1)
		go func(_usr string) {
			defer wg.Done()

			params := map[string]interface{}{
				"_user": _usr,
			}

			_user := User{
				address: _usr,
				votes: big.NewInt(0),
			}

			callObj := transactions.CallBuilder(DELEGATION, "getUserICXDelegation", params)
			res, err := c.Call(callObj)
			if err != nil {
				panic(err)
			}

			resSlice, ok := res.([]interface{})
			if !ok {
				panic("err")
			}

			for _, resMap := range resSlice {
				voteMap, ok := resMap.(map[string]interface{})
				if !ok {
					panic("err")
				}

				vote := VoteInfo{
					Address:    voteMap["_address"].(string), // is address of validator
					VotesInIcx: voteMap["_votes_in_icx"].(string),
					VotesInPer: voteMap["_votes_in_per"].(string),
				}

				if vote.Address == validator {
					fmt.Println("vote on us by: ", _usr, vote.VotesInIcx)
					_amount := util.HexToBigInt(vote.VotesInIcx)
					_updatedAmount := new(big.Int).Add(_user.votes, _amount)
					_user.votes = _updatedAmount
					// _user := User{
					// 	address: _usr,
					// 	votes:   util.HexToBigInt(vote.VotesInIcx),
					// }


					mu.Lock()
					validatorVotes = append(validatorVotes, _user)
					mu.Unlock()
				}
			}
		}(user)
	}
	wg.Wait()
	fmt.Printf("GetValidatorVotes took %.2f seconds\n", time.Since(start).Seconds())
	return validatorVotes
}

// returns the total amount of icx votes on the node
func GetOMMTotalVotes(c *client.ClientV3, validator string) *big.Int {
	validatorVotes := GetValidatorVotes(c, validator)
	amount := big.NewInt(0)

	for _, user := range validatorVotes {
		amount.Add(amount, user.votes)
	}

	return amount
}

func ExportOMMVoters(c *client.ClientV3, validator string) {
	voters := GetValidatorVotes(c, validator)
	now := time.Now()
	timestamp := now.Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("reports/%s.json", timestamp)

	voterMap := make(map[string]*big.Int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, voter := range voters {
		wg.Add(1)

		go func(_voter User) {
			defer wg.Done()
			mu.Lock()
			voterMap[_voter.address] = _voter.votes
			mu.Unlock()
		}(voter)
	}

	wg.Wait()

	data, err := json.MarshalIndent(voterMap, "", "	")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
