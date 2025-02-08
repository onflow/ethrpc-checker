package main

import (
	"errors"
	"fmt"

	"github.com/onflow/ethrpc-checker/config"
	"github.com/onflow/ethrpc-checker/rpc"
	"github.com/onflow/ethrpc-checker/types"
)

// DoTests performs the tests on the given rpc endpoint and rich privkey
// and returns an error if any of the tests failed
func DoTests(
	rpcEndpoint string,
	richPrivkey string,
	timeout string,
) error {

	conf := &config.Config{
		RpcEndpoint: rpcEndpoint,
		RichPrivKey: richPrivkey,
		Timeout:     timeout,
	}

	rCtx, err := rpc.NewContext(conf)
	if err != nil {
		return fmt.Errorf("failed to create context: %v", err)
	}

	rCtx = MustLoadContractInfo(rCtx)

	// Collect json rpc results
	var results []*types.RpcResult

	rpcs := []struct {
		name types.RpcName
		test rpc.CallRPC
	}{
		{rpc.SendRawTransaction, rpc.RpcSendRawTransactionTransferValue},
		{rpc.SendRawTransaction, rpc.RpcSendRawTransactionDeployContract},
		{rpc.SendRawTransaction, rpc.RpcSendRawTransactionTransferERC20},
		{rpc.GetBlockNumber, rpc.RpcGetBlockNumber},
		{rpc.GetGasPrice, rpc.RpcGetGasPrice},
		{rpc.GetMaxPriorityFeePerGas, rpc.RpcGetMaxPriorityFeePerGas},
		{rpc.GetChainId, rpc.RpcGetChainId},
		{rpc.GetBalance, rpc.RpcGetBalance},
		{rpc.GetTransactionCount, rpc.RpcGetTransactionCount},
		{rpc.GetBlockByHash, rpc.RpcGetBlockByHash},
		{rpc.GetBlockByNumber, rpc.RpcGetBlockByNumber},
		{rpc.GetBlockReceipts, rpc.RpcGetBlockReceipts},
		{rpc.GetTransactionByHash, rpc.RpcGetTransactionByHash},
		{rpc.GetTransactionByBlockHashAndIndex, rpc.RpcGetTransactionByBlockHashAndIndex},
		{rpc.GetTransactionByBlockNumberAndIndex, rpc.RpcGetTransactionByBlockNumberAndIndex},
		{rpc.GetTransactionReceipt, rpc.RpcGetTransactionReceipt},
		{rpc.GetBlockTransactionCountByHash, rpc.RpcGetBlockTransactionCountByHash},
		{rpc.GetCode, rpc.RpcGetCode},
		{rpc.GetStorageAt, rpc.RpcGetStorageAt},
		{rpc.NewFilter, rpc.RpcNewFilter},
		{rpc.GetFilterLogs, rpc.RpcGetFilterLogs},
		{rpc.NewBlockFilter, rpc.RpcNewBlockFilter},
		{rpc.GetFilterChanges, rpc.RpcGetFilterChanges},
		{rpc.UninstallFilter, rpc.RpcUninstallFilter},
		{rpc.GetLogs, rpc.RpcGetLogs},
		{rpc.EstimateGas, rpc.RpcEstimateGas},
		{rpc.Call, rpc.RPCCall},
	}

	for _, r := range rpcs {
		_, err := r.test(rCtx)
		if err != nil {
			// add error to results
			results = append(results, &types.RpcResult{
				Method: r.name,
				Status: types.Error,
				ErrMsg: err.Error(),
			})
			continue
		}
	}
	results = append(results, rCtx.AlreadyTestedRPCs...)

	errMsg := ""
	failed := false
	for _, r := range results {
		if r.Status == types.Error {
			failed = true
			errMsg += fmt.Sprintf(`
			%s failed: %s
			`, r.Method, r.ErrMsg)
		}
	}

	if failed {
		return errors.New(errMsg)
	}

	return nil
}
