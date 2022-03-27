package store

import "github.com/zhukovaskychina/xmysql-server/server/innodb/store/pages"

type IBuf struct {
	iBufPage pages.IBufBitMapPage
}

func (b *IBuf) GetSerializeBytes() []byte {
	return b.iBufPage.GetSerializeBytes()
}

func NewIBuf(spaceId uint32) *IBuf {
	ibuf := pages.NewIBufBitMapPage(spaceId)
	return &IBuf{ibuf}
}

//用于复盘从文件中加载出来的字节流
func NewIBufByLoadBytes(content []byte) *IBuf {

	var iBufBitMapPage = new(pages.IBufBitMapPage)
	iBufBitMapPage.FileHeader = pages.NewFileHeader()
	iBufBitMapPage.FileTrailer = pages.NewFileTrailer()

	iBufBitMapPage.LoadFileHeader(content[0:38])
	iBufBitMapPage.ChangeBufferBitMap = content[38 : 38+9192]
	iBufBitMapPage.EmptySpace = content[16384-8-8146 : 16384-8]
	iBufBitMapPage.LoadFileTrailer(content[16384-8 : 16384])

	return &IBuf{iBufPage: *iBufBitMapPage}
}
