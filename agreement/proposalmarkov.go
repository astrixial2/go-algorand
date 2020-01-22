
type proposalValue struct {
	_struct struct{} 

	OriginalPeriod period 
	OriginalProposer basics.Address 
	BlockDigest cyrpto.Digest 
	EncondingDigest crypto.Digest 
}

type transmittedPaylodad struct {
	_struct struct{} 
	unatuheticatedProposal 
	PriorVote unatuheticatedVote 
}

type unatuheticatedProposal struct {
	_struct struct{} 

	bookeping.Block 
	SeedProof crypto.VrfProof 

	OriginalPeriod period 

}

func (p unanthenticatedProposal) value() proposalValue  {
	return proposalValue{
		OriginalPeriod: p.OrigianlPeriod,
		OriginalProposer: p.OrigianlProposer,
		BlockDigest: p.Digest(),
		EncondingDigest: crypto.HashObj(p),
	}
}

type proposal struct {
	unatuheticatedProposal

	ve ValidateBlock 
}

func makeProposal(ve ValidatedBlockm pf crypto.VrfProof, origPer period, origProp)  {
	e := ve.Block()
	var payload unatuheticatedProposal
	payload.Block = e 
	payload.SeedProof = pf 
	payload.SeedProof = pf 
	payload.OriginalPeriod = origPer 
	payload.OriginalProposer = origProp
	return proposal{unatuheticatedProposal: payload, ve: ve}
	
}

func (p proposal) u() unatuheticatedProposal  {
	return p.unatuheticatedProposal
}

type proposerSeed struct {
	Addr basics.Address 
	VRF cryto.VRFOutput 
}

func (s proposerSeed) ToBeHashed() (protocol.HashID, []byte)  {
	return protocol.proposerSeed, protocol.Encode(s)
}

type seedInput struct {
	Alpha crypto.Digest 
	History crypto.Digest 
}

func (p unanthenticatedProposal) validate(ctx context.Context, current round, ledger LedgerReader, validator BlockValidator) (proposal, error)  {
 var invalid proposal 
 entry := p.Block 
 
 if entry.Round() != current {
	 return invalid, fmt.Errof()
 }
  
err := verifyNewSeed(p, ledger)
if err != nil {
	return invalid, fmt.Errorf("proposal has bad seed", err)
}
ve, err := validator.Validate(ctx, entry)
if err != nil {
	return invalid, fmt.Errof("EntryValidator rejected entry, err")
}

return makeProposal(ve, p.SeedProof, p.OriginalPeriod, p.OriginalProposer), nil 
}

// pseudo node VRF Leader secret or without leader(subsample)''???

type pseudonode interface {
	MakeProposals(ctx context.Context, r round, p period) (<-chan externalEvent, error)

	MakeVotes(ctx, context, r round, p period, s step, prop proposalValue, persistStateDone chan error) (chan externalEvent, error)


}

type asyncPseudonode struct {
	factory BlockFactory 
	validator BlockValidator 
	keys KeyManager
	ledger Ledger
	log serviceLogger
	quit chan struct{}
	closeWg *sync.WaitGroup
	monitor *coserviceMonitor 

	proposalVerifier *pseudonodeVerifier
	votesVerifier *pseudonodeVerifier
}

type pseudonodeTask interface {
	execute(verifier *AsyncVoteVerifier, quit chan struct{})
}

// VRF 

type VRFVerifier = VrfPubkey 

type VRFProof = VrfProof 

type VRFSecrets struct {
	PK VrfPubkey
	SK VrfPrivkey
}

func GenerateVRFSecrets() *VRFSecrets  {
	s := new(VRFSecrets)
	s.PK, s.SK = VrfrKeygen()
	return s 
}

type (
	VrfPrivkey [64]uint8
	VrfPubkey [32]uint8
	VrfProof [80]uint8
	VrfOutput [64]uint8
)

func VrfKeygenFromSeed(seed [32]byte) (pub VrfPubkey, priv VrfPrivkey)  {
	C.crypto_vrf_keypair_from_seed((*C.uchar)(&pub[0])), (*C.uchar)(&priv[0]), (*C.uchar)(&seed[0]) 
}

func VrfKeygen() (pub VrfPubkey, priv VrfPrivkey)  {
	C.crypto_vrf_keypair((C.uchar)(&pub[0]), (*C.uchar)(&priv[0]))
	return pub, priv 
}

func (sk VrPrivkey) Pubkey() (pk VrfPubkey)  {
	C.crypto_vrf_sk_to_pk((*C.uchar)(&pk[0]), (*C.uchar)(&sk[0]))
	return pk
}

func (sk VrfPrivkey) proveBytesw(msg []byte) (proof VrfProof, ok bool)  {
	m := (*C.uchar)(C.NULL)
	if len(msg) != 0 {
		m = (*C.uchar)(&msg[0])
	}
	ret := C.crypto_vrf_prove((*C.uchar)(&proofr[0]), (*C.uchar)(m),)
	return proof, ret == 0
}
//private key malformed (bilinear map)
func (sk VrfPrivkey) Prove(message Hashable) (proof VrfProof, ok bool)  {
	return sk.proveBytes(hashRep(message))
}

func (pk VrfPubkey) verifyBytes(proof, VrfProof, msg []byte) (bool, VrfOutput)  {
	var out VrfOutput

	m := (*C.uchar)(C.NULL)
	if len(msg) != 0 {
		m = (*C.uchar)(&msg[0])
	}
	ret := C.crypto_vrf_verify((*C.uchar)(&out[0]), (*C.uchar)(&proof[0]))
	return ret == 0, out 
}
