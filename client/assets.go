package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"io/ioutil"

	"encoding/hex"

	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
)

type DenomUnits struct {
	Denom    string `mapstructure:"denom" json:"denom"`
	Aliases    []string `mapstructure:"aliases" json:"aliases"`
	
}

type Asset struct {
	Base    string `mapstructure:"base" json:"base"`
	DenomUnits []DenomUnits `mapstructure:"denom_units" json:"denom_units"`
}

type AssetList struct {
	ChainName string  `mapstructure:"chain_name" json:"chain_name"`
	Assets    []Asset `mapstructure:"assets" json:"assets"`
}

var IbcDenomRegex = `ibc\/[0-9a-fA-F]{64}`

func ReplaceIbcWithBaseDenom(ctx Context,str string) (string, error) {
	regex := regexp.MustCompile(IbcDenomRegex)
	assetList, err := getAssetList(ctx)
	if err != nil {
		return "", err
	}
	matches := regex.FindStringSubmatch(str)
	for _, match := range matches {
		displayDenom, err := getDisplayDenom(ctx, match, assetList)
		if err != nil {
			return "", err
		}
		str = strings.Replace(str, match, displayDenom, 1)
	}
	return str, nil
}

func getAssetList(ctx Context) ([]Asset, error) {
	url := ctx.AssetListUrl
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
    }
	data, err := ioutil.ReadAll(resp.Body)
	var assetList AssetList
	err = json.Unmarshal(data, &assetList)
	if err != nil {
		return nil, err
	}
	return assetList.Assets, nil
}

func getDisplayDenom(ctx Context, ibcDenom string, assetList []Asset) (string, error) {
	// Validate ibc denom first
	err := validateIBCDenom(ibcDenom)
	if err != nil {
		return "", err
	}
	for _, asset := range assetList {
		if asset.Base == ibcDenom {
			return asset.DenomUnits[0].Aliases[0], nil
		}
	}
	return "", fmt.Errorf("Can not find denom %s in asset list", ibcDenom)
}

func validateIBCDenom(denom string) error {
	if err := sdk.ValidateDenom(denom); err != nil {
		return err
	}

	denomSplit := strings.SplitN(denom, "/", 2)

	switch {
	case denom == "ibc":
		return fmt.Errorf("denomination should be prefixed with the format 'ibc/{hash(trace + \"/\" + %s)}'", denom)

	case len(denomSplit) == 2 && denomSplit[0] == "ibc":
		if strings.TrimSpace(denomSplit[1]) == "" {
			return fmt.Errorf("denomination should be prefixed with the format 'ibc/{hash(trace + \"/\" + %s)}'", denom)
		}

		if _, err := parseHexHash(denomSplit[1]); err != nil {
			return fmt.Errorf("invalid denom trace hash %s", denomSplit[1])
		}
	}

	return nil
}

// ParseHexHash parses a hex hash in string format to bytes and validates its correctness.
func parseHexHash(hexHash string) (tmbytes.HexBytes, error) {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, err
	}

	if err := tmtypes.ValidateHash(hash); err != nil {
		return nil, err
	}

	return hash, nil
}
