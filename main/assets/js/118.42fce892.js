(window.webpackJsonp=window.webpackJsonp||[]).push([[118],{577:function(g,C,I){"use strict";I.r(C);var s=I(1),t=Object(s.a)({},(function(){var g=this,C=g.$createElement,I=g._self._c||C;return I("ContentSlotsDistributor",{attrs:{"slot-key":g.$parent.slotKey}},[I("h1",{attrs:{id:"위임자-가이드라인-cli"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#위임자-가이드라인-cli"}},[g._v("#")]),g._v(" 위임자 가이드라인 (CLI)")]),g._v(" "),I("p",[g._v("이 문서는 위임자가 커맨드라인 인터페이스(CLI, Command-Line Interface)를 통해 코스모스 허브와 소통하기 위해 필요한 모든 정보를 포함하고 있습니다.")]),g._v(" "),I("p",[g._v("또한 계정 관리, 코스모스 펀드레이저로 받은 계정을 복구하는 방법, 그리고 렛저 나노 하드웨어 지갑 사용법 또한 포함되어있습니다.")]),g._v(" "),I("p",[I("strong",[g._v("중요")]),g._v(": 이 문서에 설명되어있는 모든 단계를 신중하게 진행하십시오. 특정 행동의 실수는 소유하고 있는 아톰의 손실을 초래할 수 있습니다. 진행 전 이 문서에 있는 모든 절차를 자세히 확인하시고 필요시 코스모스 팀에게 연락하십시오. "),I("strong",[g._v("이 문서는 참고용 정보를 제공하기 위해 번역된 영어 원문의 번역본입니다. 이 문서에 포함되어있는 정보의 완결성은 보장되지 않으며, 개인의 행동에 따른 손실을 책임지지 않습니다. 꼭 영어 원문을 참고하시기 바랍니다. 만약 이 문서의 정보와 영어 원문의 정보가 다른 경우, 영어 문서의 정보가 상위 권한을 가지게 됩니다.")])]),g._v(" "),I("p",[g._v("CLI를 사용하는 위임자는 매우 실험적인 블록체인 기술이 사용되고 있는 코스모스 허브를 사용하게됩니다. 코스모스 허브는 우수한 기술을 기반으로 다수의 보안 감사를 진행했으나 문제, 업데이트 그리고 버그가 존재할 수 있습니다. 또한 블록체인 기술을 사용하는 것은 상당한 기술적 배경을 필요로 하며, 공식 팀의 컨트롤 밖에 있는 리스크가 따릅니다. 유저는 이 소프트웨어를 사용함으로써 암호학 기반 소프트웨어를 사용하는 리스크를 인지하고 있음을 인정하는 것입니다. (참고 문서: "),I("a",{attrs:{href:"https://github.com/cosmos/cosmos/blob/master/fundraiser/Interchain%20Cosmos%20Contribution%20Terms%20-%20FINAL.pdf",target:"_blank",rel:"noopener noreferrer"}},[g._v("인터체인 코스모스 펀드레이저 약관"),I("OutboundLink")],1),g._v(").")]),g._v(" "),I("p",[g._v("인터체인 재단(Interchain Foundation)과 텐더민트 팀은 소프트웨어 사용으로 발생하는 모든 손실에 대해서 책임을 지지 않습니다. Apache 2.0 라이선스 기반의 오픈소스 소프트웨어를 사용하는 것은 각 개인의 책임이며, 소프트웨어는 그 어떤 보증과 조건이 없는 'As Is(있는 그대로)' 기반으로 제공됩니다.")]),g._v(" "),I("p",[g._v("모든 절차는 신중하게 진행하시기 바랍니다.")]),g._v(" "),I("h2",{attrs:{id:"목차"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#목차"}},[g._v("#")]),g._v(" 목차")]),g._v(" "),I("ul",[I("li",[I("a",{attrs:{href:"#installing-gaiad"}},[I("code",[g._v("gaiad")]),g._v(" 설치하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#cosmos-accounts"}},[g._v("코스모스 계정")]),g._v(" "),I("ul",[I("li",[I("a",{attrs:{href:"#restoring-an-account-from-the-fundraiser"}},[g._v("펀드레이저 계정 복구하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#creating-an-account"}},[g._v("계정 생성하기")])])])]),g._v(" "),I("li",[I("a",{attrs:{href:"#accessing-the-cosmos-hub-network"}},[g._v("코스모스 허브 네트워크 액세스하기")]),g._v(" "),I("ul",[I("li",[I("a",{attrs:{href:"#running-your-own-full-node"}},[g._v("자체 풀노드 운영하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#connecting-to-a-remote-full-node"}},[g._v("원격 풀노드 연결하기")])])])]),g._v(" "),I("li",[I("a",{attrs:{href:"#setting-up-gaiad"}},[I("code",[g._v("gaiad")]),g._v(" 설정하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#querying-the-state"}},[g._v("상태(state) 조회하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#sending-transactions"}},[g._v("트랜잭션 전송하기")]),g._v(" "),I("ul",[I("li",[I("a",{attrs:{href:"#a-note-on-gas-and-fees"}},[g._v("가스(Gas)와 수수료 관련 정보")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#bonding-atoms-and-withdrawing-rewards"}},[g._v("아톰 위임 및 위임 보상 수령하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#participating-in-governance"}},[g._v("거버넌스에 참여하기")])]),g._v(" "),I("li",[I("a",{attrs:{href:"#signing-transactions-from-an-offline-computer"}},[g._v("오프라인 컴퓨터에서 트랜잭션 서명하기")])])])])]),g._v(" "),I("h2",{attrs:{id:"gaiad-설치하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#gaiad-설치하기"}},[g._v("#")]),g._v(" "),I("code",[g._v("gaiad")]),g._v(" 설치하기")]),g._v(" "),I("p",[I("code",[g._v("gaiad")]),g._v(": "),I("code",[g._v("gaiad")]),g._v("는 "),I("code",[g._v("gaiad")]),g._v(" 풀노드와 소통하기 위해 사용되는 명령어 기반 인터페이스입니다.")]),g._v(" "),I("p",[g._v("::: 경고\n"),I("strong",[g._v("추가적인 행동을 진행하기 전 최신 "),I("code",[g._v("gaiad")]),g._v(" 클라이언트를 다운로드 하셨는지 확인하십시오")]),g._v("\n:::")]),g._v(" "),I("p",[g._v("["),I("strong",[g._v("바이너리 설치하기")]),g._v("]")]),g._v(" "),I("p",[I("a",{attrs:{href:"https://cosmos.network/docs/gaia/installation.html",target:"_blank",rel:"noopener noreferrer"}},[I("strong",[g._v("소스에서 설치하기")]),I("OutboundLink")],1)]),g._v(" "),I("div",{staticClass:"custom-block tip"},[I("p",{staticClass:"custom-block-title"},[g._v("`gaiad`는 터미널 환경에서 사용됩니다. 다음과 같이 터미널을 실행하세요:")]),g._v(" "),I("ul",[I("li",[I("strong",[g._v("윈도우")]),g._v(": "),I("code",[g._v("시작")]),g._v(" > "),I("code",[g._v("모든 프로그램")]),g._v(" > "),I("code",[g._v("악세서리")]),g._v(" > "),I("code",[g._v("명령 프롬프트")])]),g._v(" "),I("li",[I("strong",[g._v("맥OS")]),g._v(": "),I("code",[g._v("파인더")]),g._v(" > "),I("code",[g._v("애플리케이션")]),g._v(" > "),I("code",[g._v("유틸리티")]),g._v(" > "),I("code",[g._v("터미널")])]),g._v(" "),I("li",[I("strong",[g._v("리눅스")]),g._v(": "),I("code",[g._v("Ctrl")]),g._v(" + "),I("code",[g._v("Alt")]),g._v(" + "),I("code",[g._v("T")]),g._v(":::")])]),g._v(" "),I("h2",{attrs:{id:"코스모스-계정"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#코스모스-계정"}},[g._v("#")]),g._v(" 코스모스 계정")]),g._v(" "),I("p",[g._v("모든 코스모스 계정에는 12개 또는 24개의 단어로 이루어진 '시드(Seed)'가 할당됩니다. 이 시드 단어(또는 시드 키)를 기반으로 다수의 코스모스 계정을 생성할 수 있습니다 (예를들어: 다수의 프라이빗 키/퍼블릭 키 쌍). 이런 형태의 월렛은 HD(Hierarchical deterministic) 월렛이라고 불립니다 (HD 월렛에 대한 자세한 정보는 "),I("a",{attrs:{href:"https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki",target:"_blank",rel:"noopener noreferrer"}},[g._v("BIP32"),I("OutboundLink")],1),g._v("를 참고하세요).")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"ICAgICAgICDqs4TsoJUgMCAgICAgICAgICAgICAgICAgICAgICAgICAgICAg6rOE7KCVIDEgICAgICAgICAgICAgICAgICAgICAgICAgICAgICDqs4TsoJUgMgoKKy0tLS0tLS0tLS0tLS0tLS0tLSsgICAgICAgICAgICAgICstLS0tLS0tLS0tLS0tLS0tLS0rICAgICAgICAgICAgICAgKy0tLS0tLS0tLS0tLS0tLS0tLSsKfCAgICAgICAgICAgICAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAgICAgICAgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICAgICAgICAgICAgIHwKfCAgICAgICDso7zshowgMCAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAg7KO87IaMIDEgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICDso7zshowgMiAgICAgIHwKfCAgICAgICAgXiAgICAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAgIF4gICAgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICAgXiAgICAgICAgIHwKfCAgICAgICAgfCAgICAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAgIHwgICAgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICAgfCAgICAgICAgIHwKfCAgICAgICAgfCAgICAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAgIHwgICAgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICAgfCAgICAgICAgIHwKfCAgICAgICAgfCAgICAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAgIHwgICAgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICAgfCAgICAgICAgIHwKfCAgICAgICAgKyAgICAgICAgIHwgICAgICAgICAgICAgIHwgICAgICAgICsgICAgICAgICB8ICAgICAgICAgICAgICAgfCAgICAgICAgKyAgICAgICAgIHwKfCAgICAg7Y2867iU66atIO2CpCAwICAgIHwgICAgICAgICAgICAgIHwgICAgIO2NvOu4lOumrSDtgqQgMSAgICAgfCAgICAgICAgICAgICAgIHwgICAgIO2NvOu4lOumrSDtgqQgMiAgICB8CnwgICAgICAgIF4gICAgICAgICB8ICAgICAgICAgICAgICB8ICAgICAgICBeICAgICAgICAgfCAgICAgICAgICAgICAgIHwgICAgICAgIF4gICAgICAgICB8CnwgICAgICAgIHwgICAgICAgICB8ICAgICAgICAgICAgICB8ICAgICAgICB8ICAgICAgICAgfCAgICAgICAgICAgICAgIHwgICAgICAgIHwgICAgICAgICB8CnwgICAgICAgIHwgICAgICAgICB8ICAgICAgICAgICAgICB8ICAgICAgICB8ICAgICAgICAgfCAgICAgICAgICAgICAgIHwgICAgICAgIHwgICAgICAgICB8CnwgICAgICAgIHwgICAgICAgICB8ICAgICAgICAgICAgICB8ICAgICAgICB8ICAgICAgICAgfCAgICAgICAgICAgICAgIHwgICAgICAgIHwgICAgICAgICB8CnwgICAgICAgICsgICAgICAgICB8ICAgICAgICAgICAgICB8ICAgICAgICArICAgICAgICAgfCAgICAgICAgICAgICAgIHwgICAgICAgICsgICAgICAgICB8CnwgICAg7ZSE65287J2067mXIO2CpCAwICAgIHwgICAgICAgICAgICAgIHwgICAgIO2UhOudvOydtOu5lyDtgqQgMSAgIHwgICAgICAgICAgICAgICB8ICAgICDtlITrnbzsnbTruZcg7YKkIDIgICB8CnwgICAgICAgIF4gICAgICAgICB8ICAgICAgICAgICAgICB8ICAgICAgICBeICAgICAgICAgfCAgICAgICAgICAgICAgIHwgICAgICAgIF4gICAgICAgICB8CistLS0tLS0tLS0tLS0tLS0tLS0rICAgICAgICAgICAgICArLS0tLS0tLS0tLS0tLS0tLS0tKyAgICAgICAgICAgICAgICstLS0tLS0tLS0tLS0tLS0tLS0rCiAgICAgICAgIHwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB8ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHwKICAgICAgICAgfCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgfAogICAgICAgICB8ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgfCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB8CiAgICAgICAgICstLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLSsKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHwKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHwKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgKy0tLS0tLS0tLSstLS0tLS0tLS0rCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIHwgICAgICAgICAgICAgICAgICAgfAogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICB8ICAgTW5lbW9uaWMgKOyLnOuTnCkgIHwKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgfCAgICAgICAgICAgICAgICAgICB8CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICstLS0tLS0tLS0tLS0tLS0tLS0tKwo="}}),g._v(" "),I("p",[g._v("특정 계좌에 보관된 자산은 프라이빗 키에 의해 관리됩니다. 이 프라이빗 키는 시드의 일방적 기능(one-way function)을 통해 생성됩니다. 프라이빗 키를 분실한 경우, 시드 키를 사용하여 프라이빗 키를 다시 복구하는 것이 가능합니다. 하지만 시드 키를 분실한 경우, 모든 프라이빗 키에 대한 사용권을 잃게 됩니다. 누군가 본인의 시드 키를 가진 경우, 해당 키와 연관된 모든 계정의 소유권을 가진 것과 동일합니다.")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("12 단어 시드키를 분실하거나 그 누구와도 공유하지 마세요. 자금 탈취와 손실을 예방하기 위해서는 다수의 시드키 사본을 만드시고 금고 같이 본인만이 알 수 있는 안전한 곳에 보관하는 것을 추천합니다. 누군가 시드키를 가지게 된 경우, 관련 프라이빗 키와 모든 계정의 소유권을 가지게 됩니다.")])])])],1),g._v(" "),I("p",[g._v("주소는 특정 계정을 구분하는 용도로 사용되며, 단어로 이루어진 특정 프리픽스(예, cosmos10)와 스트링 값을 조합한 값입니다 (예, "),I("code",[g._v("cosmos10snjt8dmpr5my0h76xj48ty80uzwhraqalu4eg")]),g._v("). 주소는 누군가 자산을 특정 계정으로 전송할때 사용되며, 퍼블릭키를 사용해 프라이빗 키를 추출하는 것은 불가능합니다.")]),g._v(" "),I("h3",{attrs:{id:"펀드레이저-계정-복구하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#펀드레이저-계정-복구하기"}},[g._v("#")]),g._v(" 펀드레이저 계정 복구하기")]),g._v(" "),I("div",{staticClass:"custom-block tip"},[I("p",[I("em",[g._v("참고: 이 항목은 코스모스 펀드레이저 참가자만을 위한 정보입니다")])])]),g._v(" "),I("p",[g._v("코스모스 펀드레이저에 참가한 인원은 12개의 단어로 구성된 시드키를 부여받습니다. 새로 생성된 시드키는 24개 단어로 이루어졌으나, 12개 단어로 이루어진 시드키 또한 모든 코스모스가 제공하는 도구에서 호환됩니다.")]),g._v(" "),I("h4",{attrs:{id:"렛저-ledger-기기-사용하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#렛저-ledger-기기-사용하기"}},[g._v("#")]),g._v(" 렛저(Ledger) 기기 사용하기")]),g._v(" "),I("p",[g._v("모든 렛저 기기에는 (코스모스 허브를 포함한) 다수의 블록체인에서 계정을 생성하기 위해 사용되는 시드키가 있습니다. 통상 시드키는 렛저 기기를 처음 활성화 할때 생성하지만, 유저가 시드키를 직접 입력하는 것 또한 가능합니다. 이제 펀드레이저를 통해 받은 시드키를 어떻게 렛저 하드웨어 지갑에 입력하는지 알아보겠습니다.")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("em",[g._v("참고: 이번 단계를 진행하실때 "),I("strong",[g._v("신규 기기를 사용하는 것을 권장합니다")]),g._v(". 한 렛저 기기에는 하나의 시드키만을 입력할 수 있습니다. 만약 이미 사용하시던 하드웨어 지갑을 사용하시기를 바라는 경우, "),I("code",[g._v("Settings")]),g._v(">"),I("code",[g._v("Device")]),g._v(">"),I("code",[g._v("Reset All")]),g._v("를 통해 리셋을 진행한 후 펀드레이저 시드를 입력할 수 있습니다. "),I("strong",[g._v("렛저 기기를 리셋할 경우, 기존에 사용했던 시드키는 기기에서 삭제됩니다. 리셋을 진행하기 전 기존 기기의 시드키를 백업하셨는지 확인하신 후 진행하시기 바랍니다.")]),g._v(" 백업 되지 않은 상태로 기기를 리셋하는 경우, 관련 계정의 자산을 잃을 수 있습니다.")])])]),g._v(" "),I("p",[g._v("다음 단계는 신규 렛저 기기 또는 초기화 된 렛저 기기에서 진행되어야 합니다:")]),g._v(" "),I("ol",[I("li",[g._v("USB를 사용해 렛저 기기를 컴퓨터에 연결하세요")]),g._v(" "),I("li",[g._v("두개의 버튼을 동시에 누르세요")]),g._v(" "),I("li",[g._v('"Restore Configuration"을 선택하세요. '),I("strong",[g._v('"Config as a new device"를 선택하시면 안됩니다')])]),g._v(" "),I("li",[g._v("원하시는 핀 번호를 입력하세요")]),g._v(" "),I("li",[g._v("12-words 옵션을 선택하세요")]),g._v(" "),I("li",[g._v("코스모스 펀드레이저에서 부여 받은 시드키를 차례대로 정확하게 입력하세요")])]),g._v(" "),I("p",[g._v("이제 렛저 하드웨어 지갑 기기가 펀드레이저 시드로 활성화되었습니다. 기존의 펀드레이저 시드를 파기하지 마십시오! 만약 렛저 기기가 고장나거나 분실된 경우, 동일한 시드키를 이용해 복구가 가능합니다.")]),g._v(" "),I("p",[g._v("이제 "),I("a",{attrs:{href:"#using-a-ledger-device"}},[g._v("여기")]),g._v("를 클릭하여 계정을 생성하는 방법을 확인하세요.")]),g._v(" "),I("h4",{attrs:{id:"컴퓨터-사용하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#컴퓨터-사용하기"}},[g._v("#")]),g._v(" 컴퓨터 사용하기")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("참고: 다음 행동은 오프라인 상태인 컴퓨터에서 진행하는 것이 더욱 안전합니다.")])])]),g._v(" "),I("p",[g._v("컴퓨터를 이용해 펀드레이저 시드키를 복구하시고 컴퓨터에 프라이빗 키를 저장사기 위해서는 다음 명령어를 실행하세요:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBrZXlzIGFkZCAmbHQ77YKkX+uqhey5rShZb3VyS2V5TmFtZSkmZ3Q7IC0tcmVjb3Zlcgo="}}),g._v(" "),I("p",[g._v("명령어를 입력하셨다면 프로그램이 지금 생성(복구)하시는 계정의 프라이빗 키를 암호화할때 사용될 비밀번호를 입력할 것을 요청합니다. 해당 계정을 이용해 트랜잭션을 보낼때마다 이 비밀번호를 입력하셔야 합니다. 만약 비밀번호를 잃어버리셨다면 시드키를 사용해 계정을 다시 복구할 수 있습니다.")]),g._v(" "),I("ul",[I("li",[I("code",[g._v("<yourKeyName>")]),g._v(" 은 계정의 이름입니다. 이는 시드키로부터 키 페어를 파생할때 레퍼런스로 사용됩니다. 이 이름은 토큰을 전송할때 보내는 계정을 구분하기 위해서 사용됩니다.")]),g._v(" "),I("li",[g._v("추가적인 선택 사항으로 명령어에 "),I("code",[g._v("--account")]),g._v(" 플래그를 추가해 특정 패스("),I("code",[g._v("0")]),g._v(", "),I("code",[g._v("1")]),g._v(", "),I("code",[g._v("2")]),g._v(", 등)를 지정할 수 있습니다. 기본적으로 "),I("code",[g._v("0")]),g._v("을 사용하여 계정이 생성됩니다.")])]),g._v(" "),I("h3",{attrs:{id:"계정-생성하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#계정-생성하기"}},[g._v("#")]),g._v(" 계정 생성하기")]),g._v(" "),I("p",[g._v("새로운 계정을 생성하기 위해서는 "),I("code",[g._v("gaiad")]),g._v("를 설치해야합니다. 신규 계정을 생성하기 전, 프라이빗 키를 어디에 저장하고 어떻게 불러올지 미리 인지를 하셔야 합니다. 프라이빗 키를 보관하기 가장 좋은 곳은 오프라인 컴퓨터 또는 렛저 하드웨어 월렛 기기입니다. 흔히 사용되는 온라인 컴퓨터에 프라이빗 키를 보관하게 될 경우, 인터넷을 통해 컴퓨터를 침투한 공격자가 프라이빗 키를 탈취할 수 있기 때문에 상당한 리스크가 존재합니다.")]),g._v(" "),I("h4",{attrs:{id:"렛저-ledger-하드웨어-월렛-기기-사용하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#렛저-ledger-하드웨어-월렛-기기-사용하기"}},[g._v("#")]),g._v(" 렛저(Ledger) 하드웨어 월렛 기기 사용하기")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("새로 주문한 렛저 기기 또는 신뢰할 수 있는 렛저 기기만을 사용하세요")])])]),g._v(" "),I("p",[g._v("렛저 기기를 처음 활성화할때 24개 단어로 구성된 시드키가 생성되고 기기에 저장됩니다. 렛저 기기의 시드키는 코스모스와 코스모스 계정과 호환이 되며, 해당 시드키를 기반으로 계정을 생성할 수 있습니다. 렛저 기기는 "),I("code",[g._v("gaiad")]),g._v("와 호환될 수 있게 설정이 되어야 합니다. 렛저 기기를 설정하는 방법은 다음과 같습니다:")]),g._v(" "),I("ol",[I("li",[I("a",{attrs:{href:"https://www.ledger.com/pages/ledger-live",target:"_blank",rel:"noopener noreferrer"}},[g._v("Ledger Live 앱"),I("OutboundLink")],1),g._v(" 다운로드")]),g._v(" "),I("li",[g._v("렛저 기기를 USB로 연결한 후 최신 펌웨어 버전으로 업데이트")]),g._v(" "),I("li",[g._v('Ledger Live 앱스토어로 이동한 후, "Cosmos" 애플리케이션 다운로드. (이 단계는 다소 시간이 걸릴 수 있습니다)')]),g._v(" "),I("li",[g._v("렛저 기기에서 코스모스 앱 선택")])]),g._v(" "),I("p",[g._v("계정을 생성하기 위해서는 다음 명령어를 실행하십시오:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBrZXlzIGFkZCAmbHQ77YKkX+uqhey5rSh5b3VyS2V5TmFtZSkmZ3Q7IC0tbGVkZ2VyIAo="}}),g._v(" "),I("ul",[I("li",[I("code",[g._v("<yourKeyName>")]),g._v(" 은 계정의 이름입니다. 이는 시드키로부터 키 페어를 파생할때 레퍼런스로 사용됩니다. 이 이름은 토큰을 전송할때 보내는 계정을 구분하기 위해서 사용됩니다.")]),g._v(" "),I("li",[g._v("추가적인 선택 사항으로 명령어에 "),I("code",[g._v("--account")]),g._v(" 플래그를 추가해 특정 패스("),I("code",[g._v("0")]),g._v(", "),I("code",[g._v("1")]),g._v(", "),I("code",[g._v("2")]),g._v(", 등)를 지정할 수 있습니다. 기본적으로 "),I("code",[g._v("0")]),g._v("을 사용하여 계정이 생성됩니다.")])]),g._v(" "),I("h4",{attrs:{id:"컴퓨터-사용하기-2"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#컴퓨터-사용하기-2"}},[g._v("#")]),g._v(" 컴퓨터 사용하기")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("참고: 다음 행동은 오프라인 상태인 컴퓨터에서 진행하는 것이 더욱 안전합니다.")])])]),g._v(" "),I("p",[g._v("계정을 생성하기 위해서는 다음 명령어를 입력하세요:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBrZXlzIGFkZCAmbHQ77YKkX+uqhey5rSh5b3VyS2V5TmFtZSkmZ3Q7Cg=="}}),g._v(" "),I("p",[g._v("위 명령어는 새로운 24단어로 구성된 시드키를 생성하고, 계정 "),I("code",[g._v("0")]),g._v("의 프라이빗 키와 퍼블릭 키를 저장합니다. 이후, 디스크에 저장될 계정 "),I("code",[g._v("0")]),g._v("의 프라이빗 키를 암호화할때 사용될 비밀번호를 입력할 것을 요청합니다. 해당 계정을 이용해 트랜잭션을 보낼때마다 이 비밀번호를 입력하셔야 합니다. 만약 비밀번호를 잃어버리셨다면 시드키를 사용해 계정을 다시 복구할 수 있습니다.")]),g._v(" "),I("div",{staticClass:"custom-block danger"},[I("p",[I("strong",[g._v("경고: 12 단어 시드키를 분실하거나 그 누구와도 공유하지 마세요. 자금 탈취와 손실을 예방하기 위해서는 다수의 시드키 사본을 만드시고 금고 같이 본인만이 알 수 있는 안전한 곳에 보관하는 것을 추천합니다. 누군가 시드키를 가지게 된 경우, 관련 프라이빗 키와 모든 계정의 소유권을 가지게 됩니다.")])])]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[g._v("시드키를 안전하게 보관하셨다면 (두번 세번씩이라도 정확하게 작성되었는지 확인하셔야 합니다!) 커맨드 라인의 기록을 다음과 같이 삭제하시면 됩니다:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"aGlzdG9yeSAtYwpybSB+Ly5iYXNoX2hpc3RvcnkK"}})],1),g._v(" "),I("ul",[I("li",[I("code",[g._v("<yourKeyName>")]),g._v(" 은 계정의 이름입니다. 이는 시드키로부터 키 페어를 파생할때 레퍼런스로 사용됩니다. 이 이름은 토큰을 전송할때 보내는 계정을 구분하기 위해서 사용됩니다.")]),g._v(" "),I("li",[g._v("추가적인 선택 사항으로 명령어에 "),I("code",[g._v("--account")]),g._v(" 플래그를 추가해 특정 패스("),I("code",[g._v("0")]),g._v(", "),I("code",[g._v("1")]),g._v(", "),I("code",[g._v("2")]),g._v(", 등)를 지정할 수 있습니다. 기본적으로 "),I("code",[g._v("0")]),g._v("을 사용하여 계정이 생성됩니다.")])]),g._v(" "),I("p",[g._v("동일한 시드키로 추가적인 계정을 생성하기 원한다면, 다음 명령어를 사용하세요:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBrZXlzIGFkZCAmbHQ77YKkX+uqhey5rSh5b3VyS2V5TmFtZSkmZ3Q7IC0tcmVjb3ZlciAtLWFjY291bnQgMQo="}}),g._v(" "),I("p",[g._v("해당 명령어는 비밀번호와 시드키를 입력할 것을 요청할 것입니다. 이 외에 추가적인 계정을 생성하시기 원한다면 account 플래그의 번호를 바꾸십시오.")]),g._v(" "),I("h2",{attrs:{id:"코스모스-허브-네트워크-사용하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#코스모스-허브-네트워크-사용하기"}},[g._v("#")]),g._v(" 코스모스 허브 네트워크 사용하기")]),g._v(" "),I("p",[g._v("블록체인의 상태(state)를 확인하거나 트랜잭션을 전송하기 위해서는 직접 풀노드를 운영하거나 다른 사람이 운영하는 풀노드에 연결할 수 있습니다.")]),g._v(" "),I("div",{staticClass:"custom-block danger"},[I("p",[I("strong",[g._v("경고: 12개 단어 / 24개 단어 시드키를 그 누구와도 공유하지 마세요. 시드키는 본인만이 알고있어야 합니다. 특히 이메일, 메시지 등의 수단으로 블록체인 서비스 지원을 사칭해 시드키를 요청할 수 있으니 주의를 바랍니다. 코스모스 팀, 텐더민트 팀 그리고 인터체인 재단은 절대로 이메일을 통해 개인 정보 또는 시드키를 요청하지 않습니다.")]),g._v(".")])]),g._v(" "),I("h3",{attrs:{id:"직접-풀노드-운영하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#직접-풀노드-운영하기"}},[g._v("#")]),g._v(" 직접 풀노드 운영하기")]),g._v(" "),I("p",[g._v("이 방법이 가장 안전한 방법이지만, 대량의 리소스를 필요로 합니다. 풀노드를 직접 운영하기 위해서는 우수한 인터넷 대역폭과 최소 1TB 상당의 하드디스크 용량을 필요로 합니다.")]),g._v(" "),I("p",[I("a",{attrs:{href:"https://cosmos.network/docs/gaia/join-mainnet.html",target:"_blank",rel:"noopener noreferrer"}},[g._v("풀노드를 운영하는 절차"),I("OutboundLink")],1),g._v("와 "),I("a",{attrs:{href:"https://cosmos.network/docs/gaia/installation.html",target:"_blank",rel:"noopener noreferrer"}},[I("code",[g._v("gaiad")]),g._v("를 설치하는 방법"),I("OutboundLink")],1),g._v("은 첨부된 링크를 확인하세요.")]),g._v(" "),I("h3",{attrs:{id:"원격-풀노드-연결하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#원격-풀노드-연결하기"}},[g._v("#")]),g._v(" 원격 풀노드 연결하기")]),g._v(" "),I("p",[g._v("만약 본인이 직접 풀노드를 운영하는 것을 원하지 않는다면 다른 사람의 풀노드에 연결을 할 수 있습니다. 이 과정에서는 신뢰할 수 있는 풀노드 운영자에만 연결하세요. 악의적인 풀노드 운영자는 트랜잭션을 막거나 틀린 정보를 전달할 가능성이 있습니다. 하지만 프라이빗 키는 당신의 컴퓨터/렛저 기기에 저장되어 있기 때문에 풀노드 운영자는 절대로 자금을 탈취할 수 없습니다. 검증된 검증인, 월렛 제공자, 거래소 등의 풀노드에만 연결하는 것을 추천드립니다.")]),g._v(" "),I("p",[g._v("풀노드에 연결하기 위해서는 다음과 같은 형식의 주소가 필요합니다: "),I("code",[g._v("https://77.87.106.33:26657")]),g._v(" ("),I("em",[g._v("이는 예시를 위한 주소이며 실제 풀노드 주소가 아닙니다")]),g._v("). 이 계정은 신뢰할 수 있는 풀노드 운영자에게서 직접 받으시기 바랍니다. 이 주소는 "),I("a",{attrs:{href:"#setting-up-gaiad"}},[g._v("다음 항목")]),g._v("에서 사용됩니다.")]),g._v(" "),I("h2",{attrs:{id:"gaiad-설정하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#gaiad-설정하기"}},[g._v("#")]),g._v(" "),I("code",[g._v("gaiad")]),g._v(" 설정하기")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[I("code",[g._v("gaiad")]),g._v("의 최신 스테이블 버전을 사용하고 있는지 확인해주세요")])])]),g._v(" "),I("p",[I("code",[g._v("gaiad")]),g._v("는 코스모스 허브 네트워크에서 운영되고 있는 노드와 소통할 수 있게 하는 도구입니다. 풀노드는 본인이 직접 운영하거나, 타인이 운영하는 풀노드를 사용할 수 있습니다. 이제 "),I("code",[g._v("gaiad")]),g._v("의 설정을 진행하겠습니다.")]),g._v(" "),I("p",[I("code",[g._v("gaiad")]),g._v("을 설정하기 위해서는 다음 명령어를 실행하세요:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBjb25maWcgJmx0O+2UjOuemOq3uChmbGFnKSZndDsgJmx0O+qwkih2YWx1ZSkmZ3Q7Cg=="}}),g._v(" "),I("p",[g._v("해당 명령어는 각 플래그에 대한 값을 설정할 수 있게 합니다. 우선 연결하고 싶은 풀노드의 주소를 입력하겠습니다:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBjb25maWcgbm9kZSAmbHQ77Zi47Iqk7Yq4KGhvc3QpJmd0OzombHQ77Y+s7Yq4KHBvcnQpJmd0OwoKLy8g7JiI7IucOiBwc3Rha2VkIGNvbmZpZyBub2RlIGh0dHBzOi8vNzcuODcuMTA2LjMzOjI2NjU3Cg=="}}),g._v(" "),I("p",[g._v("만약 풀노드를 직접 운영하시는 경우, "),I("code",[g._v("tcp://localhost:26657")]),g._v("을 주소 값으로 입력하세요.")]),g._v(" "),I("p",[g._v("이제 "),I("code",[g._v("--trust-node")]),g._v(" 플래그의 값을 설정하겠습니다:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBjb25maWcgdHJ1c3Qtbm9kZSBmYWxzZQoKLy8g66eM7JW9IOudvOydtO2KuCDtgbTrnbzsnbTslrjtirgg64W465Oc66W8IOyatOyYge2VmOqzoCDsi7bsnLzsi6Ag6rK97JqwIGB0cnVlYCDqsJLsnYQg7J6F66Cl7ZWY7IS47JqULiDqt7jroIfsp4Ag7JWK7J2AIOqyveyasCBgZmFsc2Vg66W8IOyeheugpe2VmOyEuOyalAo="}}),g._v(" "),I("p",[g._v("마지막으로 소통하고 싶은 블록체인의 "),I("code",[g._v("chain-id")]),g._v("를 입력하겠습니다:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBjb25maWcgY2hhaW4taWQgZ29zLTMK"}}),g._v(" "),I("h2",{attrs:{id:"블록체인-상태-조회하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#블록체인-상태-조회하기"}},[g._v("#")]),g._v(" 블록체인 상태 조회하기")]),g._v(" "),I("p",[I("a",{attrs:{href:"https://cosmos.network/docs/gaia/gaiad.html",target:"_blank",rel:"noopener noreferrer"}},[I("code",[g._v("gaiad")]),I("OutboundLink")],1),g._v("는 계정 잔고, 스테이킹 중인 토큰 수량, 지급 가능한 보상, 거버넌스 프로포절 등 블록체인과 관련된 모든 정보를 확인할 수 있게 합니다. 다음은 위임자에게 유용한 명령어들입니다. 다음 명령어를 실행하기 전 "),I("a",{attrs:{href:"#setting-up-gaiad"}},[g._v("gaiad 설정")]),g._v("을 진행하세요.")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"Ly8g6rOE7KCVIOyelOqzoOyZgCDqs4TsoJUg6rSA66CoIOygleuztCDsobDtmowKcHN0YWtlZCBxdWVyeSBhY2NvdW50CgovLyDqsoDspp3snbgg66qp66GdIOyhsO2ajApwc3Rha2VkIHF1ZXJ5IHZhbGlkYXRvcnMKCi8vIOqygOymneyduCDso7zshozroZwgKOyYiOyLnDogY29zbW9zMTBzbmp0OGRtcHI1bXkwaDc2eGo0OHR5ODB1endocmFxYWx1NGVnKSDqsoDspp3snbgg7KCV67O0IOyhsO2ajApwc3Rha2VkIHF1ZXJ5IHZhbGlkYXRvciAmbHQ76rKA7Kad7J24X+yjvOyGjCh2YWxpZGF0b3JBZGRyZXNzKSZndDsKCi8vIOychOyehOyekCDso7zshozroZwgKOyYiOyLnDogY29zbW9zMTBzbmp0OGRtcHI1bXkwaDc2eGo0OHR5ODB1endocmFxYWx1NGVnKSDqs4TsoJXsnZgg66qo65OgIOychOyehCDquLDroZ0g7KGw7ZqMCnBzdGFrZWQgcXVlcnkgZGVsZWdhdGlvbnMgJmx0O+ychOyehOyekF/so7zshowoZGVsZWdhdG9yQWRkcmVzcykmZ3Q7CgovLyDsnITsnoTsnpDqsIAg7Yq57KCVIOqygOymneyduOyXkOqyjCDsnITsnoTtlZwg6riw66GdIOyhsO2ajApwc3Rha2VkIHF1ZXJ5IGRlbGVnYXRpb25zICZsdDvsnITsnoTsnpBf7KO87IaMKGRlbGVnYXRvckFkZHJlc3MpJmd0OyAmbHQ76rKA7Kad7J24X+yjvOyGjCh2YWxpZGF0b3JBZGRyZXNzKSZndDsKCi8vIOychOyehOyekCDso7zshozroZwgKOyYiOyLnDogY29zbW9zMTBzbmp0OGRtcHI1bXkwaDc2eGo0OHR5ODB1endocmFxYWx1NGVnKSDsnITsnoTsnpAg66as7JuM65OcIOyhsO2ajApwc3Rha2VkIHF1ZXJ5IGRpc3RyaWJ1dGlvbiByZXdhcmRzICZsdDvsnITsnoTsnpBf7KO87IaMKGRlbGVnYXRvckFkZHJlc3MpJmd0OyAKCi8vIOuztOymneq4iChkZXBvc2l0KeydhCDrjIDquLDspJHsnbgg66qo65OgIO2UhOuhnO2PrOygiCDsobDtmowKcHN0YWtlZCBxdWVyeSBwcm9wb3NhbHMgLS1zdGF0dXMgZGVwb3NpdF9wZXJpb2QKCi8vIO2IrO2RnOqwgCDqsIDriqXtlZwg66qo65OgIO2UhOuhnO2PrOygiCDsobDtmowKcHN0YWtlZCBxdWVyeSBwcm9wb3NhbHMgLS1zdGF0dXMgdm90aW5nX3BlcmlvZAoKLy8g7Yq57KCVIO2UhOuhnO2PrOygiCBJROuhnCDtlITroZztj6zsoIgg7KCV67O0IOyhsO2ajApwc3Rha2VkIHF1ZXJ5IHByb3Bvc2FsICZsdDtwcm9wb3NhbElEJmd0Owo="}}),g._v(" "),I("p",[g._v("더 많은 명령어를 확인하기 위해서는 다음 명령어를 실행하세요:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBxdWVyeQo="}}),g._v(" "),I("p",[g._v("각 명령어에는 "),I("code",[g._v("-h")]),g._v(" 또는 "),I("code",[g._v("--help")]),g._v("를 추가하여 관련 정보를 확인하실 수 있습니다.")]),g._v(" "),I("h2",{attrs:{id:"트랜잭션-전송하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#트랜잭션-전송하기"}},[g._v("#")]),g._v(" 트랜잭션 전송하기")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[g._v("코스모스 메인넷에서는 "),I("code",[g._v("uatom")]),g._v(" 단위가 표준 단위로 사용됩니다. "),I("code",[g._v("1atom = 1,000,000uatom")]),g._v("으로 환산됩니다.")])]),g._v(" "),I("h3",{attrs:{id:"가스와-수수료에-대해서"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#가스와-수수료에-대해서"}},[g._v("#")]),g._v(" 가스와 수수료에 대해서")]),g._v(" "),I("p",[g._v("코스모스 허브 네트워크는 트랜잭션 처리를 위해 트랜잭션 수수료를 부과합니다. 해당 수수료는 트랜잭션을 실행하기 위한 가스로 사용됩니다. 공식은 다음과 같습니다:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"",base64:"7IiY7IiY66OMKEZlZSkgPSDqsIDsiqQoR2FzKSAqIOqwgOyKpCDqsJIoR2FzUHJpY2VzKQo="}}),g._v(" "),I("p",[g._v("위 공식에서 "),I("code",[g._v("gas")]),g._v("는 전송하는 트랜잭션에 따라 다릅니다. 다른 형태의 트랜잭션은 각자 다른 "),I("code",[g._v("gas")]),g._v("량을 필요로 합니다. "),I("code",[g._v("gas")]),g._v(" 수량은 트랜잭션이 실행될때 계산됨으로 사전에 정확한 값을 확인할 수 있는 방법은 없습니다. 다만, "),I("code",[g._v("gas")]),g._v(" 플래그의 값을 "),I("code",[g._v("auto")]),g._v("로 설정함으로 예상 값을 추출할 수는 있습니다. 예상 값을 수정하기 위해서는 "),I("code",[g._v("--gas-adjustment")]),g._v(" (기본 값 "),I("code",[g._v("1.0")]),g._v(") 플래그 값을 변경하셔서 트랜잭션이 충분한 가스를 확보할 수 있도록 하십시오.")]),g._v(" "),I("p",[I("code",[g._v("gasPrice")]),g._v("는 각 "),I("code",[g._v("gas")]),g._v(" 유닛의 가격입니다. 각 검증인은 직접 최소 가스 가격인 "),I("code",[g._v("min-gas-price")]),g._v("를 설정하며, 트랜잭션의 "),I("code",[g._v("gasPrice")]),g._v("가 설정한 "),I("code",[g._v("min-gas-price")]),g._v("보다 높을때 트랜잭션을 처리합니다.")]),g._v(" "),I("p",[g._v("트랜잭션 피("),I("code",[g._v("fees")]),g._v(")는 "),I("code",[g._v("gas")]),g._v(" 수량과 "),I("code",[g._v("gasPrice")]),g._v("를 곱한 값입니다. 유저는 3개의 값 중 2개의 값을 입력하게 됩니다. "),I("code",[g._v("gasPrice")]),g._v("가 높을수록 트랜잭션이 블록에 포함될 확률이 높아집니다.")]),g._v(" "),I("div",{staticClass:"custom-block tip"},[I("p",[g._v("메인넷 권장 "),I("code",[g._v("gas-prices")]),g._v("는 "),I("code",[g._v("0.0025uatom")]),g._v(" 입니다.")])]),g._v(" "),I("h3",{attrs:{id:"토큰-전송하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#토큰-전송하기"}},[g._v("#")]),g._v(" 토큰 전송하기")]),g._v(" "),I("div",{staticClass:"custom-block tip"},[I("p",{staticClass:"custom-block-title"},[g._v("**아톰을 위임하거나 위임 보상을 수령하기 전에 `gaiad`를 설치하시고 계정을 만드셔야 합니다**:::")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",{staticClass:"custom-block-title"},[g._v("참고: 다음 명령어는 온라인 상태인 컴퓨터에서 실행되어야 합니다. 해당 명령은 렛저 하드웨어 월렛 기기를 사용해 실행하는 것을 추천드립니다. 오프라인으로 트랜잭션을 발생하는 방법을 확인하기 위해서는 [여기](#signing-transactions-from-an-offline-computer)를 참고하세요")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"Ly/tirnsoJUg7IiY65+J7J2YIO2GoO2BsOydhCDsp4DsoJXtlZwg7KO87IaM66GcIOyghOyGoe2VmOq4sAovL+2MjOudvOuvuO2EsCDqsJIg7JiI7IucKOyLpOygnCDthqDtgbAg7KCE7Iah7IucIOyCrOyaqe2VmOyngCDrp4jshLjsmpQhKTogJmx0O+yImOyLoOyekF/so7zshoxf7JiI7IucJmd0Oz1jb3Ntb3MxNm05M2ZlemZpZXpodm5qYWp6cmZ5c3ptbDhxbTkyYTB3NjdudGpoZDNkMCAmbHQ77IiY65+JX+yYiOyLnCZndDs9MTAwMDAwMHVhdG9tCi8v7ZSM656Y6re4IOqwkiDsmIjsi5w6ICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDs9MC4wMDI1dWF0b20KCnBzdGFrZWQgdHggc2VuZCAmbHQ77IiY7Iug7J6QX+yjvOyGjCZndDsgJmx0O+uztOuCtOuKlF/siJjrn4kmZ3Q7IC0tZnJvbSAmbHQ77YKkX+ydtOumhCZndDsgLS1nYXMgYXV0byAtLWdhcy1hZGp1c3RtZW50IDEuNSAtLWdhcy1wcmljZXMgJmx0O+qwgOyKpF/qsIDqsqkoZ2FzUHJpY2UpJmd0Owo="}}),g._v(" "),I("h3",{attrs:{id:"아톰-위임하기-리워드-수령하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#아톰-위임하기-리워드-수령하기"}},[g._v("#")]),g._v(" 아톰 위임하기 / 리워드 수령하기")]),g._v(" "),I("div",{staticClass:"custom-block tip"},[I("p",{staticClass:"custom-block-title"},[g._v("**아톰을 위임하거나 위임 보상을 수령하기 전에 `gaiad`를 설치하시고 계정을 만드셔야 합니다**:::")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("아톰을 위임하기 전에 "),I("a",{attrs:{href:"https://cosmos.network/resources/delegators",target:"_blank",rel:"noopener noreferrer"}},[g._v("위임자 faq"),I("OutboundLink")],1),g._v("를 먼저 확인하시고 위임에 따르는 책임과 위험을 사전에 인지하시기 바랍니다")])])])])],1)]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("참고: 다음 명령어는 온라인 상태인 컴퓨터에서 실행되어야 합니다. 해당 명령은 렛저 하드웨어 월렛 기기를 사용해 실행하는 것을 추천드립니다. 오프라인으로 트랜잭션을 발생하는 방법을 확인하기 위해서는 "),I("a",{attrs:{href:"#signing-transactions-from-an-offline-computer"}},[g._v("여기")]),g._v("를 참고하세요.")])])]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"Ly8g7Yq57KCVIOqygOymneyduOyXkOqyjCDslYTthrAg7JyE7J6E7ZWY6riwIAovLyDtlIzrnpjqt7gg6rCSIOyYiOyLnDogJmx0O+qygOymneyduF/so7zshowodmFsaWRhdG9yQWRkcmVzcykmZ3Q7PSBjb3Ntb3N2YWxvcGVyMTh0aGFta2huajl3ejhwYTRuaG5wOXJsZHByZ2FudDU3cGsybThzLCAmbHQ77JyE7J6E7ZWgX+yImOufiShhbW91bnRUb0JvdW5kKSZndDs9MTAwMDAwMDB1YXRvbSwgJmx0O+qwgOyKpF/qsIDqsqkoZ2FzUHJpY2UpJmd0Oz0wLjAwMjV1YXRvbQoKcHN0YWtlZCB0eCBzdGFraW5nIGRlbGVnYXRlICZsdDvqsoDspp3snbhf7KO87IaMKHZhbGlkYXRvckFkZHJlc3MmZ3Q7ICZsdDvsnITsnoTtlaBf7IiY65+JKGFtb3VudFRvQm9uZCkmZ3Q7IC0tZnJvbSAmbHQ77JyE7J6E7J6QX+2CpF/rqoXsua0oZGVsZWdhdG9yS2V5TmFtZSkmZ3Q7IC0tZ2FzIGF1dG8gLS1nYXMtYWRqdXN0bWVudCAxLjUgLS1nYXMtcHJpY2VzICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDsKCi8vIOychOyehOuQnCDslYTthrDsnYQg64uk66W4IOqygOymneyduOyXkOqyjCDsnqzsnITsnoQg7ZWY6riwCi8vIOydtOuvuCDqsoDspp3snbjsl5Dqsowg7JyE7J6E7J20IOuQnCDsg4Htg5zsl5DshJzrp4wg7IKs7Jqp7ZWY7IukIOyImCDsnojsirXri4jri6QKLy8g7J6s7JyE7J6E7J2AIOymieyLnCDrsJjsmIHrkKnri4jri6QuIOyerOychOyehCDrjIDquLAg6riw6rCE7J2AIOyXhuyKteuLiOuLpAovLyDsnqzsnITsnoTsnYQg7KeE7ZaJ7ZWY7IugIO2bhCwg64+Z7J287ZWcIOyVhO2GsOyXkCDrjIDtlZwg7J6s7JyE7J6E7J2AIDPso7wg7ZuEIOqwgOuKpe2VqeuLiOuLpC4KLy8g7ZSM656Y6re4IOqwkiDsmIjsi5w6ICZsdDvquLDsobRf6rKA7Kad7J24X+yjvOyGjChzcmNWYWxpZGF0b3JBZGRyZXNzKSZndDs9Y29zbW9zdmFsb3BlcjE4dGhhbWtobmo5d3o4cGE0bmhucDlybGRwcmdhbnQ1N3BrMm04cywgJmx0O+yerOychOyehO2VoF/siJjrn4kmZ3Q7PTEwMDAwMDAwMHVhdG9tLCAmbHQ76rCA7IqkX+qwgOqyqShnYXNQcmljZSkmZ3Q7PTAuMDAyNXVhdG9tCgpwc3Rha2VkIHR4IHN0YWtpbmcgcmVkZWxlZ2F0ZSAmbHQ76riw7KG0X+qygOymneyduF/so7zshowoc3JjVmFsaWRhdG9yQWRkcmVzcykmZ3Q7ICZsdDvsnbTrj5ntlaBf6rKA7Kad7J24X+yjvOyGjChkZXN0VmFsaWRhdG9yQWRkcmVzcykmZ3Q7ICZsdDvsnqzsnITsnoTtlaBf7IiY65+JKGFtb3VudFRvUmVkZWxlZ2F0ZSkmZ3Q7IC0tZnJvbSAmbHQ77JyE7J6E7J6QX+2CpF/rqoXsua0oZGVsZWdhdG9yS2V5TmFtZSkmZ3Q7IC0tZ2FzIGF1dG8gLS1nYXMtYWRqdXN0bWVudCAxLjUgLS1nYXMtcHJpY2VzICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDsKCi8vIOuqqOuToCDrpqzsm4zrk5wg7IiY66C57ZWY6riwCi8vIO2UjOuemOq3uCDqsJIg7JiI7IucOiAmbHQ76rCA7IqkX+qwgOqyqShnYXNQcmljZSkmZ3Q7PTAuMDAyNXVhdG9tCgpwc3Rha2VkIHR4IGRpc3RyaWJ1dGlvbiB3aXRoZHJhdy1hbGwtcmV3YXJkcyAtLWZyb20gJmx0O+ychOyehOyekF/tgqRf66qF7LmtKGRlbGVnYXRvcktleU5hbWUpJmd0OyAtLWdhcyBhdXRvIC0tZ2FzLWFkanVzdG1lbnQgMS41IC0tZ2FzLXByaWNlcyAmbHQ76rCA7IqkX+qwgOqyqShnYXNQcmljZSkmZ3Q7CgovLyDtirnsoJUg6rKA7Kad7J247Jy866GcIOu2gO2EsCDsnITsnoQg7Leo7IaM7ZWY6riwCi8vIOychOyehCDst6jshozqsIAg7JmE66OM65CY6riwIOychO2VtOyEnOuKlCAz7KO87J2YIOq4sOqwhOydtCDqsbjrpqzrqbAsIOychOyehCDst6jshozqsIAg7KeE7ZaJ7KSR7J24IOq4sOqwhOyXkOuKlCDtlbTri7kg7JWE7Yaw7J2EIOyghOyGoe2VmOyLpCDsiJgg7JeG7Iq164uI64ukLgovLyDtlIzrnpjqt7gg6rCSIOyYiOyLnDogJmx0O+qygOymneyduF/so7zshowodmFsaWRhdG9yQWRkcmVzcykmZ3Q7PWNvc21vc3ZhbG9wZXIxOHRoYW1raG5qOXd6OHBhNG5obnA5cmxkcHJnYW50NTdwazJtOHMsICZsdDvsnITsnoRf7Leo7IaM7ZWgX+yImOufiShhbW91bnRUb1VuYm9uZCkmZ3Q7PTEwMDAwMDAwdWF0b20sICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDs9MC4wMDI1dWF0b20KCnBzdGFrZWQgdHggc3Rha2luZyB1bmJvbmQgJmx0O+qygOymneyduF/so7zshowodmFsaWRhdG9yQWRkcmVzcykmZ3Q7ICZsdDvsnITsnoRf7Leo7IaM7ZWgX+yImOufiShhbW91bnRUb1VuYm9uZCkmZ3Q7IC0tZnJvbSAmbHQ77JyE7J6E7J6QX+2CpF/rqoXsua0oZGVsZWdhdG9yS2V5TmFtZSkmZ3Q7IC0tZ2FzIGF1dG8gLS1nYXMtYWRqdXN0bWVudCAxLjUgLS1nYXMtcHJpY2VzICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDsKCjo6OiB0aXAK66Cb7KCAIOq4sOq4sOulvCDsgqzsmqntlbQg7Yq4656c7J6t7IWY7J2EIOuwnOyDne2VmOyLnOuKlCDqsr3smrAsIOugm+yggCDquLDquLDsl5DshJwg7Yq4656c7J6t7IWY7J2EIO2ZleyduO2VmOuKlCDqs7zsoJXsnbQg7LaU6rCA7KCB7Jy866GcIOuwnOyDneuQqeuLiOuLpC4g7Lu07ZOo7YSw7JeQIOyXsOqysOuQmOyWtCDsnojripQg6riw6riw7JeQ7IScIO2KuOuenOyereyFmOydhCDshJzrqoXtlZjshZTslbwg64Sk7Yq47JuM7YGs66GcIOyghO2MjOuQqeuLiOuLpC4KOjo6IAoK7ZW064u5IO2KuOuenOyereyFmOydtCDshLHqs7XsoIHsnLzroZwg7KeE7ZaJ65CcIOqyg+ydhCDtmZXsnbjtlZjquLAg7JyE7ZW07ISc64qUIOuLpOydjCDsobDtmowg66qF66C57Ja066W8IOyCrOyaqe2VmOyEuOyalDoKCmBgYGJhc2gKLy8g7JWE7Yaw7J2EIOychOyehO2VmOqxsOuCmCDrpqzsm4zrk5zrpbwg7IiY66C57ZWY7IugIO2bhCDqs4TsoJUg7J6U6rOg6rCAIOuLrOudvOynkeuLiOuLpCAo6rOE7KCVIOyelOqzoCDtmZXsnbgg66qF66C57Ja0KQpwc3Rha2VkIHF1ZXJ5IGFjY291bnQKCi8vIOychOyehOydhCDsp4TtlontlZjshajri6TrqbQg7Iqk7YWM7J207YK5IOyelOqzoOqwgCDtkZzsi5zrkKnri4jri6QgKOyKpO2FjOydtO2CuSDtmZXsnbgg66qF66C57Ja0KQpwc3Rha2VkIHF1ZXJ5IGRlbGVnYXRpb25zICZsdDvsnITsnoTsnpBf7KO87IaMKGRlbGVnYXRvckFkZHJlc3MpJmd0OwoKLy8g7Yq4656c7J6t7IWY7J20IOu4lOuhneyytOyduOyXkCDtj6ztlajrkJjsl4jsnLzrqbQg7ZW064u5IHR4IOygleuztOulvCDsoITri6ztlanri4jri6QKLy8g7Yq4656c7J6t7IWY7J2EIOyDneyEse2VmOyFqOydhOuVjCDtkZzsi5zrkJjsl4jrjZggdHggaGFzaOulvCDsnoXroKXtlZjshLjsmpQgKO2KuOuenOyereyFmCDtmZXsnbgg66qF66C57Ja0KQpwc3Rha2VkIHF1ZXJ5IHR4ICZsdDvtirjrnpzsnq3shZhf7ZW07IucKHR4SGFzaCkmZ3Q7Cgo="}}),g._v(" "),I("p",[g._v("만약 원격 풀노드를 사용해 트랜잭션을 전송하신 경우, 블록 익스플로러를 통해 트랜잭션을 확인하십시오.")]),g._v(" "),I("h2",{attrs:{id:"거버넌스-참가하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#거버넌스-참가하기"}},[g._v("#")]),g._v(" 거버넌스 참가하기")]),g._v(" "),I("h3",{attrs:{id:"거버넌스에-대해서"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#거버넌스에-대해서"}},[g._v("#")]),g._v(" 거버넌스에 대해서")]),g._v(" "),I("p",[g._v("코스모스 허브는 아톰을 스테이킹 한 위임자들이 투표를 할 수 있는 시스템이 내장되어있습니다. 프로포절의 종류는 3개가 있으며 다음과 같습니다:")]),g._v(" "),I("ul",[I("li",[I("code",[g._v("텍스트 프로포절(Text Proposals)")]),g._v(": 가장 기본적인 형태의 프로포절입니다. 특정 주제에 대한 네트워크의 의견을 확인하기 위해서 사용됩니다.")]),g._v(" "),I("li",[I("code",[g._v("파라미터 프로포절(Parameter Proposals)")]),g._v(": 네트워크의 기존 파라미터 값을 변경하는 것을 제안하기 위해서 사용됩니다.")]),g._v(" "),I("li",[I("code",[g._v("소프트웨어 업그레이드 프로포절(Software Upgrade Proposal)")]),g._v(": 코스모스 허브의 소프트웨어를 업그레이드 하는 것을 제안하기 위해서 사용됩니다.")])]),g._v(" "),I("p",[g._v("모든 아톰 보유자는 프로포절을 제안할 수 있습니다. 특정 프로포절의 투표가 활성화되기 위해서는 "),I("code",[g._v("minDeposit")]),g._v("값에 정의된 보증금 보다 높은 "),I("code",[g._v("deposit")]),g._v(" 비용이 예치되어야 합니다. "),I("code",[g._v("deposit")]),g._v("은 프로포절 제안자 외에도 보증금을 추가할 수 있습니다. 만약 제안자가 필요한 보증금 보다 낮은 보증금을 입금한 경우, 프로포절은 "),I("code",[g._v("deposit_period")]),g._v(" 상태로 들어가 추가 보증금 입금을 대기합니다. 모든 아톰 보유자는 "),I("code",[g._v("depositTx")]),g._v(" 트랜잭션을 통해 보증금을 추가할 수 있습니다.")]),g._v(" "),I("p",[g._v("프로포절의 "),I("code",[g._v("deposit")]),g._v("이 "),I("code",[g._v("minDeposit")]),g._v("을 도달하게 되면 해당 프로포절의 2주 간의 "),I("code",[g._v("voting_period")]),g._v("(투표 기간)이 시작됩니다. "),I("strong",[g._v("위임된 아톰")]),g._v("의 보유자는 해당 프로포절에 투표를 행사할 수 있으며, "),I("code",[g._v("Yes")]),g._v(", "),I("code",[g._v("No")]),g._v(", "),I("code",[g._v("NoWithVeto")]),g._v(" 또는 "),I("code",[g._v("Abstain")]),g._v(" 표를 선택할 수 있습니다. 각 표는 투표자의 위임된 아톰 수량을 반영하게 됩니다. 만약 위임자가 직접 투표를 진행하지 않은 경우, 위임자는 검증인의 표를 따르게 됩니다. 하지만 모든 위임자는 직접 투표를 행사하여 검증인의 표와 다른 표를 행사할 수 있습니다.")]),g._v(" "),I("p",[g._v("투표 기간이 끝난 후, 프로포절이 50% 이상의 "),I("code",[g._v("Yes")]),g._v("표를 받았고 ("),I("code",[g._v("Abstain")]),g._v(" 표를 제외하고) "),I("code",[g._v("NoWithVeto")]),g._v(" ("),I("code",[g._v("Abstain")]),g._v(" 표를 제외하고) 표가 33.33% 이하일 경우 통과하게 됩니다.")]),g._v(" "),I("h3",{attrs:{id:"거버넌스-참여하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#거버넌스-참여하기"}},[g._v("#")]),g._v(" 거버넌스 참여하기")]),g._v(" "),I("div",{staticClass:"custom-block warning"},[I("p",[I("strong",[g._v("참고: 다음 명령어는 온라인 상태인 컴퓨터에서만 진행이 가능합니다. 해당 명령은 렛저 하드웨어 월렛 기기를 사용해 실행하는 것을 추천드립니다. 오프라인으로 트랜잭션을 발생하는 방법을 확인하기 위해서는 "),I("a",{attrs:{href:"#signing-transactions-from-an-offline-computer"}},[g._v("여기")]),g._v("를 참고하세요.")])])]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"Ly8g7ZSE66Gc7Y+s7KCIIOygnOyViO2VmOq4sAovLyAmbHQ77ZSE66Gc7Y+s7KCIX+yiheulmCh0eXBlKSZndDs9dGV4dC9wYXJhbWV0ZXJfY2hhbmdlL3NvZnR3YXJlX3VwZ3JhZGUKLy8g7ZSM656Y6re4IOqwkiDsmIjsi5w6ICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDs9MC4wMDI1dWF0b20KCnBzdGFrZWQgdHggZ292IHN1Ym1pdC1wcm9wb3NhbCAtLXRpdGxlICZxdW90O1Rlc3QgUHJvcG9zYWwmcXVvdDsgLS1kZXNjcmlwdGlvbiAmcXVvdDtNeSBhd2Vzb21lIHByb3Bvc2FsJnF1b3Q7IC0tdHlwZSAmbHQ77ZSE66Gc7Y+s7KCIX+yiheulmCh0eXBlKSZndDsgLS1kZXBvc2l0PTEwMDAwMDAwdWF0b20gLS1nYXMgYXV0byAtLWdhcy1wcmljZXMgJmx0O+qwgOyKpF/qsIDqsqkoZ2FzUHJpY2UpJmd0OyAtLWZyb20gJmx0O+ychOyehOyekF/tgqRf66qF7LmtKGRlbGVnYXRvcktleU5hbWUpJmd0OwoKLy8g7ZSE66Gc7Y+s7KCI7J2YIOuztOymneq4iCDstpTqsIDtlZjquLAKLy8g7ZSE66Gc7Y+s7KCI7J2YIHByb3Bvc2FsSUQg7KGw7ZqMOiAkZ2FpYWQgcXVlcnkgZ292IHByb3Bvc2FscyAtLXN0YXR1cyBkZXBvc2l0X3BlcmlvZAovLyDtjIzrnbzrr7jthLAg6rCSIOyYiOyLnDogJmx0O+uztOymneq4iChkZXBvc2l0KSZndDs9MTAwMDAwMDB1YXRvbQoKcHN0YWtlZCB0eCBnb3YgZGVwb3NpdCAmbHQ77ZSE66Gc7Y+s7KCIX0lEKHByb3Bvc2FsSUQpJmd0OyAmbHQ77LaU6rCA7ZWgX+uztOymneq4iChkZXBvc2l0KSZndDsgLS1nYXMgYXV0byAtLWdhcy1wcmljZXMgJmx0O+qwgOyKpF/qsIDqsqkoZ2FzUHJpY2UpJmd0OyAtLWZyb20gJmx0O+ychOyehOyekF/tgqRf66qF7LmtKGRlbGVnYXRvcktleU5hbWUpJmd0OwoKLy8g7ZSE66Gc7Y+s7KCI7JeQIO2IrO2RnO2VmOq4sAovLyDtlITroZztj6zsoIjsnZggcHJvcG9zYWxJRCDsobDtmow6ICRnYWlhZCBxdWVyeSBnb3YgcHJvcG9zYWxzIC0tc3RhdHVzIHZvdGluZ19wZXJpb2QgCi8vICZsdDvtkZxf7ISg7YOdKG9wdGlvbikmZ3Q7PXllcy9uby9ub193aXRoX3ZldG8vYWJzdGFpbgoKcHN0YWtlZCB0eCBnb3Ygdm90ZSAmbHQ77ZSE66Gc7Y+s7KCIX0lEKHByb3Bvc2FsSUQpJmd0OyAmbHQ77ZGcX+yEoO2DnShvcHRpb24pJmd0OyAtLWdhcyBhdXRvIC0tZ2FzLXByaWNlcyAmbHQ76rCA7IqkX+qwgOqyqShnYXNQcmljZSkmZ3Q7IC0tZnJvbSAmbHQ77JyE7J6E7J6QX+2CpF/rqoXsua0oZGVsZWdhdG9yS2V5TmFtZSkmZ3Q7Cg=="}}),g._v(" "),I("h2",{attrs:{id:"오프라인-컴퓨터에서-트랜잭션-서명하기"}},[I("a",{staticClass:"header-anchor",attrs:{href:"#오프라인-컴퓨터에서-트랜잭션-서명하기"}},[g._v("#")]),g._v(" 오프라인 컴퓨터에서 트랜잭션 서명하기")]),g._v(" "),I("p",[g._v("렛저 기기가 없거나 오프라인 컴퓨터에서 프라이빗 키를 관리하고 싶으신 경우, 다음 절차를 따라하세요. 우선 "),I("strong",[g._v("온라인")]),g._v(" 컴퓨터에서 미서명 트랜잭션을 다음과 같이 생성하십시오 (다음 예시에서는 위임 트랜잭션을 예시로 사용합니다):")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"Ly8g7JWE7YawIOuzuOuUqe2VmOq4sCAKLy8g7ZSM656Y6re4IOqwkiDsmIjsi5w6ICZsdDvrs7jrlKntlaAg7IiY65+JKGFtb3VudFRvQm9uZCkmZ3Q7PTEwMDAwMDAwdWF0b20sICZsdDvsnITsnoTtlaAg6rKA7Kad7J247J2YIGJlY2gzMiDso7zshowoYmVjaDMyQWRkcmVzc09mVmFsaWRhdG9yKSZndDs9Y29zbW9zdmFsb3BlcjE4dGhhbWtobmo5d3o4cGE0bmhucDlybGRwcmdhbnQ1N3BrMm04cywgJmx0O+qwgOyKpCDqsIDqsqkoZ2FzUHJpY2UpJmd0Oz0wLjAwMjV1YXRvbQoKcHN0YWtlZCB0eCBzdGFraW5nIGRlbGVnYXRlICZsdDvqsoDspp3snbhf7KO87IaMKHZhbGlkYXRvckFkZHJlc3MpJmd0OyAmbHQ77JyE7J6E7ZWgX+yImOufiShhbW91bnRUb0JvbmQpJmd0OyAtLWZyb20gJmx0O+ychOyehOyekF/so7zshowoZGVsZWdhdG9yQWRkcmVzcykmZ3Q7IC0tZ2FzIGF1dG8gLS1nYXMtYWRqdXN0bWVudCAxLjUgLS1nYXMtcHJpY2VzICZsdDvqsIDsiqRf6rCA6rKpKGdhc1ByaWNlKSZndDsgLS1nZW5lcmF0ZS1vbmx5ICZndDsgdW5zaWduZWRUWC5qc29uCg=="}}),g._v(" "),I("p",[g._v("서명을 진행하기 위해서는 "),I("code",[g._v("chain-id")]),g._v(", "),I("code",[g._v("account-number")]),g._v(", 그리고 "),I("code",[g._v("sequence")]),g._v(" 값이 필요합니다. "),I("code",[g._v("chain-id")]),g._v("는 트랜잭션을 전송할 블록체인의 고유 식별 번호입니다. "),I("code",[g._v("account-number")]),g._v("는 계정이 처음 자산을 받을 때 생성되는 고유 번호입니다. "),I("code",[g._v("sequence")]),g._v("는 리플레이 공격을 방지하기 위해 전송한 트랜잭션의 수량을 기록하는데 사용됩니다.")]),g._v(" "),I("p",[g._v("체인 아이디(chain-id) 값은 해당 블록체인의 제네시스 파일에서 받으실 수 있습니다 (현재 기준 코스모스 허브는 "),I("code",[g._v("cosmoshub-2")]),g._v("). account-number와 sequence는 계정 조회 명령어를 사용해 확인하실 수 있습니다.")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCBxdWVyeSBhY2NvdW50ICZsdDvqs4TsoJVf7KO87IaMKHlvdXJBZGRyZXNzKSZndDsgLS1jaGFpbi1pZCBjb3Ntb3NodWItMgo="}}),g._v(" "),I("p",[g._v("이후 서명이 진행되지 않은 "),I("code",[g._v("unsignedTx.json")]),g._v(" 파일을 복사하신 후 (USB 등을 이용하여) 오프라인 컴퓨터로 이동하십시오. 만약 오프라인 컴퓨터에 아직 계정을 생성하지 않으셨을 경우, "),I("a",{attrs:{href:"#using-a-computer"}},[g._v("이 항목")]),g._v("을 참고하여 오프라인 컴퓨터에서 계정을 생성하세요. 안전을 위해서 서명하기 전에 다음 명령어를 실행해 트랜잭션의 파라미터를 한번 더 확인하십시오:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"Y2F0IHVuc2lnbmVkVHguanNvbgo="}}),g._v(" "),I("p",[g._v("이제 다음 명령어를 실행해 트랜잭션을 서명합니다:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCB0eCBzaWduIHVuc2lnbmVkVHguanNvbiAtLWZyb20gJmx0O+ychOyehOyekF/tgqRf66qF7LmtKGRlbGVnYXRvcktleU5hbWUpJmd0OyAtLW9mZmxpbmUgLS1jaGFpbi1pZCBjb3Ntb3NodWItMiAtLXNlcXVlbmNlICZsdDvsi5ztgIDsiqQoc2VxdWVuY2UpJmd0OyAtLWFjY291bnQtbnVtYmVyICZsdDvqs4TsoJVf67KI7Zi4KGFjY291bnROdW1iZXIpJmd0OyAmZ3Q7IHNpZ25lZFR4Lmpzb24K"}}),g._v(" "),I("p",[g._v("서명된 "),I("code",[g._v("signedTx.json")]),g._v(" 파일을 복사하시고 다시 온라인 컴퓨터로 이동하세요. 다음 명령어를 실행해 해당 트랜잭션을 네트워크에 전파하세요:")]),g._v(" "),I("tm-code-block",{staticClass:"codeblock",attrs:{language:"bash",base64:"cHN0YWtlZCB0eCBicm9hZGNhc3Qgc2lnbmVkVHguanNvbgo="}})],1)}),[],!1,null,null,null);C.default=t.exports}}]);