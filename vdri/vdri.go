package vdri

type VDRI interface {
	GetDIDDoc(did string) CommonDIDDoc
}

type CommonDIDDoc interface {
	GetServicePoint(serviceid string) string
}
