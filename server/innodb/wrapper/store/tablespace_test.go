package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zhukovaskychina/xmysql-server/server/conf"
	"testing"
)

func TestNewSysTableSpace(t *testing.T) {
	cfg := conf.NewCfg()
	cfg.DataDir = "/tmp/"
	NewSysTableSpace(cfg)

}

func TestNewTableSpaceFile(t *testing.T) {
	cfg := conf.NewCfg()
	cfg.DataDir = "/tmp/"
	NewTableSpaceFile(cfg, "student", "student", 1, false)

}

func TestCheckTableSpaceFileHeader(t *testing.T) {
	cfg := conf.NewCfg()
	cfg.DataDir = "/tmp/"

	ts := NewSysTableSpace(cfg)
	ts.GetFirstFsp()
	ts.GetDictTable()
}

//测试段的建立
func TestCreateSegment(t *testing.T) {
	cfg := conf.NewCfg()
	cfg.DataDir = "/tmp/"
	ts := NewTableSpaceFile(cfg, "student", "student", 1, false)
	fmt.Println(ts)
}

func TestSegmentInit(t *testing.T) {

	//seg:=segs.NewIBDSegment(1,"",false,false,false,false)

	cfg := conf.NewCfg()

	cfg.DataDir = "/tmp/"
	ts := NewTableSpaceFile(cfg, "student", "student", 1, false)
	fmt.Println(ts)

}

func Test8thData(t *testing.T) {
	cfg := conf.NewCfg()
	cfg.DataDir = "/Users/zhukovasky/xmysql/data"
	cfg.BaseDir = "/Users/zhukovasky/xmysql"
	sysTs := NewSysTableSpace(cfg)
	pageBytes, _ := sysTs.LoadPageByPageNumber(8)

	index := NewPageIndexByLoadBytesWithTuple(pageBytes, NewSysTableTuple())
	assert.Equal(t, int(index.GetPageNumber()), 8)

	assert.Equal(t, len(index.SlotRowData.FullRowList()), 3)
}
