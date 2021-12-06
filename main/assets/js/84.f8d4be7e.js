(window.webpackJsonp=window.webpackJsonp||[]).push([[84],{692:function(e,o,t){"use strict";t.r(o);var a=t(1),r=Object(a.a)({},(function(){var e=this,o=e.$createElement,t=e._self._c||o;return t("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[t("h1",{attrs:{id:"gov-subspace"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#gov-subspace"}},[e._v("#")]),e._v(" "),t("code",[e._v("gov")]),e._v(" subspace")]),e._v(" "),t("p",[e._v("The "),t("code",[e._v("gov")]),e._v(" module is responsible for on-chain governance proposals and voting functionality.")]),e._v(" "),t("table",[t("tr",[t("th",[e._v("Key")]),e._v(" "),t("th",[e._v("Value")])]),e._v(" "),e._l(e.$themeConfig.currentParameters.gov,(function(o,a){return t("tr",[t("td",[t("a",{attrs:{href:"#"+a}},[t("code",[e._v(e._s(a))])])]),e._v(" "),t("td",[t("code",[e._v(e._s(o))])])])}))],2),e._v(" "),t("p",[e._v("The "),t("code",[e._v("gov")]),e._v(" module is responsible for the on-chain governance system. In this system, holders of the native staking token of the chain may vote on proposals on a 1-token per 1-vote basis. The module supports:")]),e._v(" "),t("ul",[t("li",[t("strong",[e._v("Proposal submission")]),e._v(": Users can submit proposals with a deposit. Once the minimum deposit is reached, proposal enters voting period")]),e._v(" "),t("li",[t("strong",[e._v("Vote")]),e._v(": Participants can vote on proposals that reached MinDeposit")]),e._v(" "),t("li",[t("strong",[e._v("Inheritance and penalties")]),e._v(": Delegators inherit their validator's vote if they don't vote themselves.")]),e._v(" "),t("li",[t("strong",[e._v("Claiming deposit")]),e._v(": Users that deposited on proposals can recover their deposits if the proposal was accepted OR if the proposal never entered voting period.")])]),e._v(" "),t("h2",{attrs:{id:"governance-notes-on-parameters"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#governance-notes-on-parameters"}},[e._v("#")]),e._v(" Governance notes on parameters")]),e._v(" "),t("h3",{attrs:{id:"depositparams"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#depositparams"}},[e._v("#")]),e._v(" "),t("code",[e._v("depositparams")])]),e._v(" "),t("h4",{attrs:{id:"min-deposit"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#min-deposit"}},[e._v("#")]),e._v(" "),t("code",[e._v("min_deposit")])]),e._v(" "),t("p",[t("strong",[e._v("The minimum deposit required for a proposal to enter the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(", in micro-ATOMs.")])]),e._v(" "),t("ul",[t("li",[e._v("on-chain value: "),t("code",[e._v(e._s(e.$themeConfig.currentParameters.gov.depositparams.min_deposit))])]),e._v(" "),t("li",[t("a",{attrs:{href:"https://www.mintscan.io/cosmos/proposals/47",target:"_blank",rel:"noopener noreferrer"}},[e._v("Proposal 47"),t("OutboundLink")],1),e._v(" change: "),t("code",[e._v("64000000")]),e._v(" "),t("code",[e._v("uatom")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-4")]),e._v(" default: "),t("code",[e._v("512000000")]),e._v(" "),t("code",[e._v("uatom")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-3")]),e._v(" default: "),t("code",[e._v("512000000")]),e._v(" "),t("code",[e._v("uatom")])])]),e._v(" "),t("p",[e._v("Prior to a governance proposal entering the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(" (ie. for the proposal to be voted upon), there must be at least a minimum number of ATOMs deposited. Anyone may contribute to this deposit. Deposits of passed and failed proposals are returned to the contributors. Deposits are burned when proposals 1) "),t("a",{attrs:{href:"#maxdepositperiod"}},[e._v("expire")]),e._v(", 2) fail to reach "),t("a",{attrs:{href:"#quorum"}},[e._v("quorum")]),e._v(", or 3) are "),t("a",{attrs:{href:"#veto"}},[e._v("vetoed")]),e._v(". This parameter subkey value represents the minimum deposit required for a proposal to enter the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(" in micro-ATOMs, where "),t("code",[e._v("512000000uatom")]),e._v(" is equivalent to 512 ATOM.")]),e._v(" "),t("h5",{attrs:{id:"decreasing-the-value-of-mindeposit"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#decreasing-the-value-of-mindeposit"}},[e._v("#")]),e._v(" Decreasing the value of "),t("code",[e._v("mindeposit")])]),e._v(" "),t("p",[e._v("Decreasing the value of the "),t("code",[e._v("mindeposit")]),e._v(" subkey will enable governance proposals to enter the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(" with fewer ATOMs at risk. This will likely increase the volume of new governance proposals.")]),e._v(" "),t("h5",{attrs:{id:"increasing-the-value-of-mindeposit"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#increasing-the-value-of-mindeposit"}},[e._v("#")]),e._v(" Increasing the value of "),t("code",[e._v("mindeposit")])]),e._v(" "),t("p",[e._v("Increasing the value of the "),t("code",[e._v("mindeposit")]),e._v(" subkey will require risking a greater number of ATOMs before governance proposals may enter the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(". This will likely decrease the volume of new governance proposals.")]),e._v(" "),t("h4",{attrs:{id:"max-deposit-period"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#max-deposit-period"}},[e._v("#")]),e._v(" "),t("code",[e._v("max_deposit_period")])]),e._v(" "),t("p",[t("strong",[e._v("The maximum amount of time that a proposal can accept deposit contributions before expiring, in nanoseconds.")])]),e._v(" "),t("ul",[t("li",[e._v("on-chain value: "),t("code",[e._v(e._s(e.$themeConfig.currentParameters.gov.depositparams.max_deposit_period))])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-4")]),e._v(" default: "),t("code",[e._v("1209600000000000")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-3")]),e._v(" default: "),t("code",[e._v("1209600000000000")])])]),e._v(" "),t("p",[e._v("Prior to a governance proposal entering the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(", there must be at least a minimum number of ATOMs deposited. This parameter subkey value represents the maximum amount of time that the proposal has to reach the minimum deposit amount before expiring. The maximum amount of time that a proposal can accept deposit contributions before expiring is currently "),t("code",[e._v("1209600000000000")]),e._v(" nanoseconds or 14 days. If the proposal expires, any deposit amounts will be burned.")]),e._v(" "),t("h5",{attrs:{id:"decreasing-the-value-of-maxdepositperiod"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#decreasing-the-value-of-maxdepositperiod"}},[e._v("#")]),e._v(" Decreasing the value of "),t("code",[e._v("maxdepositperiod")])]),e._v(" "),t("p",[e._v("Decreasing the value of the "),t("code",[e._v("maxdepositperiod")]),e._v(" subkey will decrease the time for deposit contributions to governance proposals. This will likely decrease the time that some proposals remain visible and potentially decrease the likelihood that they will enter the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(". This may increase the likelihood that proposals will expire and have their deposits burned.")]),e._v(" "),t("h5",{attrs:{id:"increasing-the-value-of-maxdepositperiod"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#increasing-the-value-of-maxdepositperiod"}},[e._v("#")]),e._v(" Increasing the value of "),t("code",[e._v("maxdepositperiod")])]),e._v(" "),t("p",[e._v("Increasing the value of the "),t("code",[e._v("maxdepositperiod")]),e._v(" subkey will extend the time for deposit contributions to governance proposals. This will likely increase the time that some proposals remain visible and potentially increase the likelihood that they will enter the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(". This may decrease the likelihood that proposals will expire and have their deposits burned.")]),e._v(" "),t("h5",{attrs:{id:"notes"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#notes"}},[e._v("#")]),e._v(" Notes")]),e._v(" "),t("p",[e._v("Currently most network explorers (eg. Hubble, Big Dipper, Mintscan) give the same visibility to proposals in the deposit period as those in the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(". This means that a proposal with a small deposit (eg. 0.001 ATOM) will have the same visibility as those with a full 512 ATOM deposit in the voting period.")]),e._v(" "),t("h3",{attrs:{id:"votingparams"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#votingparams"}},[e._v("#")]),e._v(" "),t("code",[e._v("votingparams")])]),e._v(" "),t("h4",{attrs:{id:"votingperiod"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#votingperiod"}},[e._v("#")]),e._v(" "),t("code",[e._v("votingperiod")])]),e._v(" "),t("p",[t("strong",[e._v("The maximum amount of time that a proposal can accept votes before the voting period concludes, in nanoseconds.")])]),e._v(" "),t("ul",[t("li",[e._v("on-chain value: "),t("code",[e._v(e._s(e.$themeConfig.currentParameters.gov.votingparams.voting_period))])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-4")]),e._v(" default: "),t("code",[e._v("1209600000000000")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-3")]),e._v(" default: "),t("code",[e._v("1209600000000000")])])]),e._v(" "),t("p",[e._v("Once a governance proposal enters the voting period, there is a maximum period of time that may elapse before the voting period concludes. This parameter subkey value represents the maximum amount of time that the proposal has to accept votes, which is currently "),t("code",[e._v("1209600000000000")]),e._v(" nanoseconds or 14 days. If the proposal vote does not reach quorum ((ie. 40% of the network's voting power is participating) before this time, any deposit amounts will be burned and the proposal's outcome will not be considered to be valid. Voters may change their vote any number of times before the voting period ends. This voting period is currently the same for any kind of governance proposal.")]),e._v(" "),t("h5",{attrs:{id:"decreasing-the-value-of-votingperiod"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#decreasing-the-value-of-votingperiod"}},[e._v("#")]),e._v(" Decreasing the value of "),t("code",[e._v("votingperiod")])]),e._v(" "),t("p",[e._v("Decreasing the value of the "),t("code",[e._v("votingperiod")]),e._v(" subkey will decrease the time for voting on governance proposals. This will likely:")]),e._v(" "),t("ol",[t("li",[e._v("decrease the proportion of the network that participates in voting, and")]),e._v(" "),t("li",[e._v("decrease the likelihood that quorum will be reached.")])]),e._v(" "),t("h5",{attrs:{id:"increasing-the-value-of-votingperiod"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#increasing-the-value-of-votingperiod"}},[e._v("#")]),e._v(" Increasing the value of "),t("code",[e._v("votingperiod")])]),e._v(" "),t("p",[e._v("Increasing the value of the "),t("code",[e._v("votingperiod")]),e._v(" subkey will increase the time for voting on governance proposals. This may:")]),e._v(" "),t("ol",[t("li",[e._v("increase the proportion of the network that participates in voting, and")]),e._v(" "),t("li",[e._v("increase the likelihood that quorum will be reached.")])]),e._v(" "),t("h5",{attrs:{id:"notes-2"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#notes-2"}},[e._v("#")]),e._v(" Notes")]),e._v(" "),t("p",[e._v("Historically, off-chain discussions and engagement appears to be have been greater occurred during the voting period of a governance proposal than when the proposal is posted off-chain as a draft. A non-trivial amount of the voting power has voted in the second week of the voting period. Proposals 23, 19, and 13 each had approximately 80% network participation or more.")]),e._v(" "),t("h3",{attrs:{id:"tallyparams"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#tallyparams"}},[e._v("#")]),e._v(" "),t("code",[e._v("tallyparams")])]),e._v(" "),t("h4",{attrs:{id:"quorum"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#quorum"}},[e._v("#")]),e._v(" "),t("code",[e._v("quorum")])]),e._v(" "),t("p",[t("strong",[e._v("The minimum proportion of network voting power required for a governance proposal's outcome to be considered valid.")])]),e._v(" "),t("ul",[t("li",[e._v("on-chain value: "),t("code",[e._v(e._s(e.$themeConfig.currentParameters.gov.tallyparams.quorum))])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-4")]),e._v(" default: "),t("code",[e._v("0.400000000000000000")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-3")]),e._v(" default: "),t("code",[e._v("0.400000000000000000")])])]),e._v(" "),t("p",[e._v("Quorum is required for the outcome of a governance proposal vote to be considered valid and for deposit contributors to recover their deposit amounts, and this parameter subkey value represents the minimum value for quorum. Voting power, whether backing a vote of 'yes', 'abstain', 'no', or 'no-with-veto', counts toward quorum. If the proposal vote does not reach quorum (ie. 40% of the network's voting power is participating) before this time, any deposit amounts will be burned and the proposal outcome will not be considered to be valid.")]),e._v(" "),t("h5",{attrs:{id:"decreasing-the-value-of-quorum"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#decreasing-the-value-of-quorum"}},[e._v("#")]),e._v(" Decreasing the value of "),t("code",[e._v("quorum")])]),e._v(" "),t("p",[e._v("Decreasing the value of the "),t("code",[e._v("quorum")]),e._v(" subkey will enable a smaller proportion of the network to legitimize the outcome of a proposal. This increases the risk that a decision will be made with a smaller proportion of ATOM-stakers' positions being represented, while decreasing the risk that a proposal will be considered invalid. This will likely decrease the risk of a proposal's deposit being burned.")]),e._v(" "),t("h5",{attrs:{id:"increasing-the-value-of-quorum"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#increasing-the-value-of-quorum"}},[e._v("#")]),e._v(" Increasing the value of "),t("code",[e._v("quorum")])]),e._v(" "),t("p",[e._v("Increasing the value of the "),t("code",[e._v("quorum")]),e._v(" subkey will require a larger proportion of the network to legitimize the outcome of a proposal. This decreases the risk that a decision will be made with a smaller proportion of ATOM-stakers' positions being represented, while increasing the risk that a proposal will be considered invalid. This will likely increase the risk of a proposal's deposit being burned.")]),e._v(" "),t("h4",{attrs:{id:"threshold"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#threshold"}},[e._v("#")]),e._v(" "),t("code",[e._v("threshold")])]),e._v(" "),t("p",[t("strong",[e._v("The minimum proportion of participating voting power required for a governance proposal to pass.")])]),e._v(" "),t("ul",[t("li",[e._v("on-chain value: "),t("code",[e._v(e._s(e.$themeConfig.currentParameters.gov.tallyparams.threshold))])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-4")]),e._v(" default: "),t("code",[e._v("0.500000000000000000")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-3")]),e._v(" default: "),t("code",[e._v("0.500000000000000000")])])]),e._v(" "),t("p",[e._v("A simple majority 'yes' vote (ie. 50% of participating voting power) is required for a governance proposal vote to pass. Though necessary, a simple majority 'yes' vote may not be sufficient to pass a proposal in two scenarios:")]),e._v(" "),t("ol",[t("li",[e._v("Failure to reach "),t("a",{attrs:{href:"#quorum"}},[e._v("quorum")]),e._v(" of 40% network power or")]),e._v(" "),t("li",[e._v("A 'no-with-veto' vote of 33.4% of participating voting power or greater.")])]),e._v(" "),t("p",[e._v("If a governance proposal passes, deposit amounts are returned to contributors. If a text-based proposal passes, nothing is enacted automatically, but there is a social expectation that participants will co-ordinate to enact the commitments signalled in the proposal. If a parameter change proposal passes, the protocol parameter will automatically change immediately after the "),t("a",{attrs:{href:"#votingperiod"}},[e._v("voting period")]),e._v(" ends, and without the need to run new software. If a community-spend proposal passes, the Community Pool balance will decrease by the number of ATOMs indicated in the proposal and the recipient's address will increase by this same number of ATOMs immediately after the voting period ends.")]),e._v(" "),t("h5",{attrs:{id:"decreasing-the-value-of-threshold"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#decreasing-the-value-of-threshold"}},[e._v("#")]),e._v(" Decreasing the value of "),t("code",[e._v("threshold")])]),e._v(" "),t("p",[e._v("Decreasing the value of the "),t("code",[e._v("threshold")]),e._v(" subkey will decrease the proportion of voting power required to pass a proposal. This may:")]),e._v(" "),t("ol",[t("li",[e._v("increase the likelihood that a proposal will pass, and")]),e._v(" "),t("li",[e._v("increase the likelihood that a minority group will effect changes to the network.")])]),e._v(" "),t("h5",{attrs:{id:"increasing-the-value-of-threshold"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#increasing-the-value-of-threshold"}},[e._v("#")]),e._v(" Increasing the value of "),t("code",[e._v("threshold")])]),e._v(" "),t("p",[e._v("Increasing the value of the "),t("code",[e._v("threshold")]),e._v(" subkey will increase the proportion of voting power required to pass a proposal. This may:")]),e._v(" "),t("ol",[t("li",[e._v("decrease the likelihood that a proposal will pass, and")]),e._v(" "),t("li",[e._v("decrease the likelihood that a minority group will effect changes to the network.")])]),e._v(" "),t("h4",{attrs:{id:"veto-threshold"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#veto-threshold"}},[e._v("#")]),e._v(" "),t("code",[e._v("veto_threshold")])]),e._v(" "),t("p",[t("strong",[e._v("The minimum proportion of participating voting power to veto (ie. fail) a governance proposal.")])]),e._v(" "),t("ul",[t("li",[e._v("on-chain value: "),t("code",[e._v(e._s(e.$themeConfig.currentParameters.gov.tallyparams.veto_threshold))])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-4")]),e._v(" default: "),t("code",[e._v("0.334000000000000000")])]),e._v(" "),t("li",[t("code",[e._v("cosmoshub-3")]),e._v(" default: "),t("code",[e._v("0.334000000000000000")])])]),e._v(" "),t("p",[e._v("Though a simple majority 'yes' vote (ie. 50% of participating voting power) is required for a governance proposal vote to pass, a 'no-with-veto' vote of 33.4% of participating voting power or greater can override this outcome and cause the proposal to fail. This enables a minority group representing greater than 1/3 of voting power to fail a proposal that would otherwise pass.")]),e._v(" "),t("h5",{attrs:{id:"decreasing-the-value-of-veto-threshold"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#decreasing-the-value-of-veto-threshold"}},[e._v("#")]),e._v(" Decreasing the value of "),t("code",[e._v("veto_threshold")])]),e._v(" "),t("p",[e._v("Decreasing the value of the "),t("code",[e._v("veto_threshold")]),e._v(" subkey will decrease the proportion of participating voting power required to veto. This will likely:")]),e._v(" "),t("ol",[t("li",[e._v("enable a smaller minority group to prevent proposals from passing, and")]),e._v(" "),t("li",[e._v("decrease the likelihood that contentious proposals will pass.")])]),e._v(" "),t("h5",{attrs:{id:"increasing-the-value-of-veto-threshold"}},[t("a",{staticClass:"header-anchor",attrs:{href:"#increasing-the-value-of-veto-threshold"}},[e._v("#")]),e._v(" Increasing the value of "),t("code",[e._v("veto_threshold")])]),e._v(" "),t("p",[e._v("Increasing the value of the "),t("code",[e._v("veto_threshold")]),e._v(" subkey will increase the proportion of participating voting power required to veto. This will require a larger minority group to prevent proposals from passing, and will likely increase the likelihood that contentious proposals will pass.")])])}),[],!1,null,null,null);o.default=r.exports}}]);