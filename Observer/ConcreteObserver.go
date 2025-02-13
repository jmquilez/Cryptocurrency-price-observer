package Observer

type ConcreteObserver struct {
	id string
	Btc float64
	Eth float64
	Ada float64
	Btc_Ok bool
	Eth_Ok bool
	Ada_Ok bool
}

func (p *ConcreteObserver) Update(Btc, Eth, Ada float64) {
	p.Btc = Btc
	p.Eth = Eth
	p.Ada = Ada
}

func (p *ConcreteObserver) GetID() string {
	return p.id
}

func (p *ConcreteObserver) GetBtc() float64 {
	return p.Btc
}

func (p *ConcreteObserver) GetEth() float64 {
	return p.Eth
}

func (p *ConcreteObserver) GetAda() float64 {
	return p.Ada
}

func (p *ConcreteObserver) GetBtc_Ok() bool {
	return p.Btc_Ok
}

func (p *ConcreteObserver) GetEth_Ok() bool {
	return p.Eth_Ok
}

func (p *ConcreteObserver) GetAda_Ok() bool {
	return p.Ada_Ok
}
