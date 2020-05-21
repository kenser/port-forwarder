package forward

import (
	"context"
	"github.com/cloverzrg/go-portforward/logger"
	"github.com/cloverzrg/go-portforward/model/forwarddao"
	"github.com/cloverzrg/go-portforward/service/forwardermanager"
	"github.com/cloverzrg/go-portforward/web/dto"
)

func Add(ctx context.Context, req dto.AddForward) (id int, err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	id, err = forwarddao.Add(req.Network, req.ListenAddress, req.ListenPort, req.TargetAddress, req.TargetPort)
	if err != nil {
		return id, err
	}
	err = forwardermanager.StartById(ctx, id)
	if err != nil {
		return id, err
	}
	err = forwarddao.UpdateByIdMap(id, map[string]interface{}{"status": 1})
	if err != nil {
		return id, err
	}
	return id, err
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

func Start(ctx context.Context, id int) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
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

func Delete(ctx context.Context, id int) (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	err = forwardermanager.CloseById(ctx, id)
	if err == forwardermanager.ErrNotRunning {
		err = nil
	}
	if err != nil {
		return err
	}
	err = forwarddao.DeleteById(id)
	if err != nil {
		return err
	}
	return err
}

func Find(ctx context.Context, filters dto.PortForwardFilters) (res dto.ForwardList, err error) {
	list, total, err := forwarddao.FindByFilters(filters)
	var resList []dto.ForwardDetail
	for _, v := range list {
		resList = append(resList, dto.ForwardDetail{
			Id:            v.Id,
			Status:        v.Status,
			Network:       v.Network,
			ListenAddress: v.ListenAddress,
			ListenPort:    v.ListenPort,
			TargetAddress: v.TargetAddress,
			TargetPort:    v.TargetPort,
		})
	}
	res = dto.ForwardList{
		Total: total,
		List:  resList,
	}
	return res, err
}

func GetDetailById(ctx context.Context, id int) (res dto.ForwardDetail, err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	data, err := forwarddao.GetById(id)
	if err != nil {
		return res, err
	}
	res = dto.ForwardDetail{
		Id:            data.Id,
		Status:        data.Status,
		Network:       data.Network,
		ListenAddress: data.ListenAddress,
		ListenPort:    data.ListenPort,
		TargetAddress: data.TargetAddress,
		TargetPort:    data.TargetPort,
	}
	return res, err
}

func StartUp() (err error) {
	defer func() {
		if err != nil {
			logger.Error(err)
		}
	}()
	list, err := forwarddao.FindAllRunning()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, v := range list {
		err = Start(ctx, v.Id)
		if err != nil {
			return err
		}
	}
	return err
}