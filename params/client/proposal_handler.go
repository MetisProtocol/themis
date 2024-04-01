package client

import (
	govclient "github.com/metis-seq/themis/gov/client"
	"github.com/metis-seq/themis/params/client/cli"
	"github.com/metis-seq/themis/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
