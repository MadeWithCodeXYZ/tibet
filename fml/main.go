package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
)

type Coin struct {
	Amount uint64 `json:"amount"`
	ParentCoinInfo string `json:"parent_coin_info"`
	PuzzleHash string `json:"puzzle_hash"`
}

type CoinSpend struct {
	Coin Coin `json:"coin"`
	PuzzleReveal string `json:"puzzle_reveal"`
	Solution string `json:"solution"`
}

type SpendBundle struct {
	AggregatedSignature string `json:"aggregated_signature"`
	CoinSpends []CoinSpend `json:"coin_spends"`
}

type MempoolItem struct {
	SpendBundle SpendBundle `json:"spend_bundle"`
}

type AllMempoolItemsResponse struct {
	Success bool `json:"success"`
	MempoolItems map[string]MempoolItem `json:"mempool_items"`
}

type GetMempoolItemByParentCoinInfoArgs struct {
    ParentCoinInfo string `json:"parent_coin_info"`
	RequestURL string `json:"request_url"`
}

func GetAllMempoolItemsResponse(request_url string) (AllMempoolItemsResponse, error) {
	res, err := http.Get(request_url)
	if err != nil {
		return AllMempoolItemsResponse{}, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AllMempoolItemsResponse{}, err
	}

	var resp AllMempoolItemsResponse
	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return AllMempoolItemsResponse{}, err
	}

	return resp, nil
}

func GetMempoolItemByParentCoinInfo(c *fiber.Ctx) error {
	args := new(GetMempoolItemByParentCoinInfoArgs)

	if err := c.BodyParser(args); err != nil {
		return err
	}

	response, err := GetAllMempoolItemsResponse(args.RequestURL)
	if err != nil {
		return c.JSON(fiber.Map{
			"item": nil,
		})
	}

	var item SpendBundle
	var found bool = false

	for _, v := range response.MempoolItems {
		for _, cs := range v.SpendBundle.CoinSpends {
			if cs.Coin.ParentCoinInfo == args.ParentCoinInfo {
				found = true
				item = v.SpendBundle
				break
			}
		}

		if found {
			break
		}
	}

	if !found {
		return c.JSON(fiber.Map{
			"item": nil,
		})
	}

	return c.JSON(fiber.Map{
		"item": item,
	})
}

func main() {
    app := fiber.New(fiber.Config{
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
    })

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Fast Mempool Locator is running! ~ FML")
    })
	app.Post("/get_mempool_item_by_parent_coin_info", GetMempoolItemByParentCoinInfo)

    app.Listen(":1337")
}