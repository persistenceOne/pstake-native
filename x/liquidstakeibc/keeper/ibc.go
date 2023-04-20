package keeper

func (k Keeper) OnChanOpenInit() (string, error) {
	return "", nil
}

func (k Keeper) OnChanOpenAck() error {
	return nil
}
func (k Keeper) OnAcknowledgementPacket() error {
	return nil
}
func (k Keeper) OnTimeoutPacket() error {
	return nil
}
