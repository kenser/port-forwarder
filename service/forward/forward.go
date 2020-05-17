package forward

import (
	"context"
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/model/forwarddao"
	"github.com/cloverzrg/go-portforward/service/forwardermanager"
	"github.com/cloverzrg/go-portforward/web/dto"
)

func Add(ctx context.Context, req dto.AddForward) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	id, err := forwarddao.Add(req.Network, req.ListenAddress, req.ListenPort, req.TargetAddress, req.TargetPort)
	if err != nil {
		return err
	}
	err = forwardermanager.StartById(ctx, id)
	if err != nil {
		return err
	}
	err = forwarddao.UpdateByIdMap(id, map[string]interface{}{"status": 1})
	if err != nil {
		return err
	}
	return err
}

func Stop(ctx context.Context, id int) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	err = forwardermanager.CloseById(ctx, id)
	if err != nil {
		return err
	}
	err = forwarddao.UpdateByIdMap(id, map[string]interface{}{"status": 0})
	if err != nil {
		return err
	}
	return err
}

func Find(ctx context.Context, filters string) (err error) {
	return err
}
