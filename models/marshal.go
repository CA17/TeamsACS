package models

import (
	"encoding/json"
	"time"

	"github.com/ca17/teamsacs/common/timeutil"
)

func (d NetCpe) MarshalJSON() ([]byte, error) {
	type Alias NetCpe
	return json.Marshal(&struct {
		Alias
		CwmpLastInform string `json:"cwmp_last_inform"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
	}{
		Alias:          (Alias)(d),
		CwmpLastInform: d.CwmpLastInform.Format(time.RFC3339),
		CreatedAt:      d.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      d.UpdatedAt.Format(time.RFC3339),
	})
}

func (d SysOprLog) MarshalJSON() ([]byte, error) {
	type Alias SysOprLog
	return json.Marshal(&struct {
		Alias
		OptTime string `json:"opt_time"`
	}{
		Alias:   (Alias)(d),
		OptTime: d.OptTime.Format(time.RFC3339),
	})
}

func (d *CwmpPresetTask) MarshalJSON() ([]byte, error) {
	type Alias CwmpPresetTask
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		ExecTime  string `json:"exec_time"`
		RespTime  string `json:"resp_time"`
	}{
		Alias:     (*Alias)(d),
		CreatedAt: d.CreatedAt.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
		UpdatedAt: d.UpdatedAt.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
		ExecTime:  d.ExecTime.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
		RespTime:  d.RespTime.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
	})
}

func (d *CwmpPresetTask) UnmarshalJSON(b []byte) error {
	type Alias CwmpPresetTask
	var tmp = struct {
		*Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		ExecTime  string `json:"exec_time"`
		RespTime  string `json:"resp_time"`
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	d.CreatedAt, _ = time.Parse(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.CreatedAt)
	d.UpdatedAt, _ = time.Parse(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.UpdatedAt)
	d.ExecTime, _ = time.Parse(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.ExecTime)
	d.RespTime, _ = time.Parse(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.RespTime)
	return nil
}
