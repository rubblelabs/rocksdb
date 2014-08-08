rocksdb
=======

Storage layer to access rippled nodestore and a tool for dumping its contents by walking the ledger from a  specified start and end ledger hash.


##Installation

Install [Go](http://golang.org/doc/install) making sure to set a GOPATH and add GOPATH/bin to your PATH.
Install [RocksDB dependencies](https://github.com/facebook/rocksdb/blob/master/INSTALL.md).

```bash
git clone https://github.com/facebook/rocksdb.git
cd rocksdb
make shared_lib
CGO_CFLAGS="-I/path/to/rocksdb/include" CGO_LDFLAGS="-L/path/to/rocksdb" go get -u -v github.com/rubblelabs/rocksdb/rdb
```

changing the /path/to/rocksdb as appropriate.


##Usage

For OSX

```bash
export DYLD_LIBRARY_PATH=/path/to/rocksdb/
rdb -help
```

For Linux

```bash
export LD_LIBRARY_PATH=/path/to/rocksdb/
rdb -help
```

##Examples

See rdb/dump.sh for an example of custom formatting and composing useful outputs from multiple commands that could be applied to map reduce scenarios.

All commands have a -start and -end flag which is a ledger hash. For example to get a summary of ledgers 100,000 to 32,570 (the default), you can do:

```bash
rdb -start=491E88B0A5AB29378B4F4E6EAB1E782AF495D712A817C943D0D7A36045EFA611 -end=4109C6F2045FC7EFF4CDE8F9905D19C28820D86304080FF886B299F0206E42B5 -path=ripple/nodedb/ -command=summary
100000,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99999,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99998,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99997,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99996,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
...
```

Show the summary of all transactions and account state nodes in a ledger, with columns:

Ledger Sequence, Transaction Inner Nodes, Payments, AccountSets, SetRegularKeys, OfferCreates, OfferCancels, TrustSets, Amendments, SetFees, Account Inner Nodes, AccountRoot, DirectoryNode, Amendments, LedgerHashes, Offer, RippleState, FeeSettings

```bash
rdb -path=ripple/nodedb/ -command=summary
100000,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99999,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99998,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99997,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
99996,0,0,0,0,0,0,0,0,0,337,213,140,0,3,30,115,0
....
```

Show the diff between one ledger's account state and its predecessor's account state, with columns:

Ledger, Deleted/Added/moved, Node Type, Depth, NodeId, NodeValue

```bash
rdb -path=ripple/nodedb/ -command=diff 
100000,D,Account Node,0,87696621BDF55B732FC417BFA3FC36D71A14A9A7DB46A8B39A2A065F58666639,0000000000000000034D494E00044B598464FCB5610BAA175D7D509330056ED7728B59D562BF8DABDBB0954B99237BCA7747F81E2EFA6CE328EDF48473CE1C5615EBDB0D41377C51E04C2C9A74620B96FCAE1E7BA4C0D5C7E6C04A3DFC67E3A91149AB6AC4A5B54D4E42444F0D6F56B010055BE886F531A674ED544446659AC433D27E1689915A2A35B51DAE3C4980EA45C8307D855918381F759012AF96B2C48596758A12CE882061FE0C6921B5C2BEAA1B67C134DE1F094ADF82B45F3E8EC74538934A92F284060E6DA160EDDAE0896A3E1BF665BF96210E6EA51A7E61131670DC8CDBE52E742119C4FF2C354E5B79213CC3CD1147742BB13B71DCF0B996EF48C16EFEC2C93A402AB17076E461E62CD882AE8F2AD9A4A32C5B78A5835605C80D995CEF99C06DA8DD784B7CBE33C28F07FB653E9EA067AD59650AB899D09BDB8A7FF86BAE16D99CCC04D6D63F0A98B039B23D6386A907DD80CF7C16611032FB3406DFE4FFA65DF42728FF8B081FF75ADD40F93FD4D77CC009AAE7A275FC1E75015269315989C5D72B829252EACB6D757386BDE712314F7776277A1C9AC25BC506B4BF66AC20E6F9E3CD64F7338306969BA157A24FE48C73C6AF3CB3CAF76433CBBAE6C525AA92D72C15239E1D443FF353AD32F468386D34D058A46BE517A7D02282BD080C29C992635F2AE093E7732FE937DE8A3524391BDF4CE9AC747CF2D060CDDEFD4958473B31E028970B
100000,A,Account Node,0,1082031456E3CE683C94368E067CF44D22B4456A793235163B333E69F4E716EA,0000000000000000034D494E00044B598464FCB5610BAA175D7D509330056ED7728B59D562BF8DABDBB0954B99237BCA7747F81E2EFA6CE328EDF48473CE1C5615EBDB0D41377C51E04C2C9A74620B96FCAE1E7BA4C0D5C7E6C04A3DFC67E3A91149AB6AC4A5B54D4E42444F0D6F56B010055BE886F531A674ED544446659AC433D27E1689915A2A35B51DAE3C4980EA45C8307D855918381F759012AF96B2C48596758A12CE882061FE0C6921B5C2BEAA1B67C134DE1F094ADF82B45F3E8EC74538934A92F284060E6DA160EDDAE0896A3E1BF665BF96210E6EA51A7E61131670DC8CDBE52E742119C4FF2C354E5B79213CC3CD1147742BB13B71DCF0B996EF48C16EFEC2C93A402AB17076E461E62CD882AE8F2AD9A4A32C5B78A5835605C80D995CEF99C06DA8DD784B7CBE33C28F07FB653E9EA067AD59650AB899D09BDB8A7FF86BAE16D99CCC04D6D63F0A98B039B23D6386A907DD80CF7C16611032FB3406DFE4FFA65DF42728FF8B089FD199A4A5187EB1DE2E2BDFC6D6C6540C1E711B38C3BAF06F0899D0BAE68A39CB6D757386BDE712314F7776277A1C9AC25BC506B4BF66AC20E6F9E3CD64F7338306969BA157A24FE48C73C6AF3CB3CAF76433CBBAE6C525AA92D72C15239E1D443FF353AD32F468386D34D058A46BE517A7D02282BD080C29C992635F2AE093E7732FE937DE8A3524391BDF4CE9AC747CF2D060CDDEFD4958473B31E028970B
100000,D,Account Node,1,1FF75ADD40F93FD4D77CC009AAE7A275FC1E75015269315989C5D72B829252EA,0000000000000000034D494E00E13D9F0A5DF5CD0BBA84937F9B4225E68EB63F1D7BFB124493C6D9EC211AB4FF92347BE0C83BB02C5683094ED6E5F90A4C6258262988811A0F3745DEDCB370914F3F130C72E9C2AADC538CF6AFDE61FEA84D0F8F563A27F03AAA45EE9F3B7DE9B7C8C4466B3FC5A21E289E98A08A7B6AD5C83AE0E4A623D011C30048666133F1028D4882965C99E04CE30B874A814048FFA875FD99865F2E809B0DA62EB635909D2ADB38DDC210988A2C8369E2792015C491B71A583B3B87450A7E0FC97117A94F96CBFDADB4FB48D1DB4B4D40DCA3187619BCAECE4B8F3C8E8D775A6A19F88285A086EA692B0953C58A6B64AEAFB9FD42F0FDC028E688E13FD1183C9BA4FC9988CDF3C2FE4657CC0F32EAD14C2D940B679A764C0E89203EF20CD36C5598240FF91300627DA9533C6CE5C4923A9E31635108DCA1628788C60914DCD9B5A603A0621E80B851CD59CA2012EA2778E345766C3D4B4BCA24486F39F7D15B956834254F417FF06DD68DAFDE51A9180674AC29862A942095487AE51C858D53F84B1BA7C01EE1B8B60DF42F0F69DFA0C79C0DFCD4D4CCDAB31B0207D357E6E4D9F33E2152BE290D8ECAABAF462E370AA7D4699F21132F3DB89FAE708C567283E408007CD6BB2B3AB8E4486FD43AADBBD54120A6D4ADF99676EF3547D5B17F43AB4594FF0000000000000000000000000000000000000000000000000000000000000000
...
```

Dump the ledger headers:

```bash
rdb -path=ripple/nodedb/ -command=ledgers 
491E88B0A5AB29378B4F4E6EAB1E782AF495D712A817C943D0D7A36045EFA611:000186A0000186A0014C575200000186A0016345785D89D9044570B5228629561CAA77851EFA880D355C6559150C97F5BB619A549DA579CF1D00000000000000000000000000000000000000000000000000000000000000001082031456E3CE683C94368E067CF44D22B4456A793235163B333E69F4E716EA188A1728188A173C0A00
4570B5228629561CAA77851EFA880D355C6559150C97F5BB619A549DA579CF1D:0001869F0001869F014C5752000001869F016345785D89D904D711E24F58292814F1782E1B978010E3035EB15A0E6C2C21D3C81A6D06DC49A1000000000000000000000000000000000000000000000000000000000000000087696621BDF55B732FC417BFA3FC36D71A14A9A7DB46A8B39A2A065F58666639188A1714188A17280A00
D711E24F58292814F1782E1B978010E3035EB15A0E6C2C21D3C81A6D06DC49A1:0001869E0001869E014C5752000001869E016345785D89D9042ACFD78893111663AE5FFBEA683FED073A15935F887E9A1E3001B907B1299E3F00000000000000000000000000000000000000000000000000000000000000001E987534476DD69123E3779C21E2E93CD6AC5938B58B88572F7227E61400BB87188A1700188A17140A00
2ACFD78893111663AE5FFBEA683FED073A15935F887E9A1E3001B907B1299E3F:0001869D0001869D014C5752000001869D016345785D89D904DB714BD4503A4096AB7BFAB5174D8D30D4FB5CF05476CF4D8F9E8077226A18B7000000000000000000000000000000000000000000000000000000000000000088CBCDA8D06B275494BD812C3DD8FB414F81300D97FB5943DE2738DA7215B413188A16EC188A17000A00
...
```

Dump the whole account state tree with columns:

Ledger, Node Type, Depth, NodeId, NodeValue

```bash
rdb -path=ripple/nodedb/ -command=dump
100000,Account Node,0,1082031456E3CE683C94368E067CF44D22B4456A793235163B333E69F4E716EA,0000000000000000034D494E00044B598464FCB5610BAA175D7D509330056ED7728B59D562BF8DABDBB0954B99237BCA7747F81E2EFA6CE328EDF48473CE1C5615EBDB0D41377C51E04C2C9A74620B96FCAE1E7BA4C0D5C7E6C04A3DFC67E3A91149AB6AC4A5B54D4E42444F0D6F56B010055BE886F531A674ED544446659AC433D27E1689915A2A35B51DAE3C4980EA45C8307D855918381F759012AF96B2C48596758A12CE882061FE0C6921B5C2BEAA1B67C134DE1F094ADF82B45F3E8EC74538934A92F284060E6DA160EDDAE0896A3E1BF665BF96210E6EA51A7E61131670DC8CDBE52E742119C4FF2C354E5B79213CC3CD1147742BB13B71DCF0B996EF48C16EFEC2C93A402AB17076E461E62CD882AE8F2AD9A4A32C5B78A5835605C80D995CEF99C06DA8DD784B7CBE33C28F07FB653E9EA067AD59650AB899D09BDB8A7FF86BAE16D99CCC04D6D63F0A98B039B23D6386A907DD80CF7C16611032FB3406DFE4FFA65DF42728FF8B089FD199A4A5187EB1DE2E2BDFC6D6C6540C1E711B38C3BAF06F0899D0BAE68A39CB6D757386BDE712314F7776277A1C9AC25BC506B4BF66AC20E6F9E3CD64F7338306969BA157A24FE48C73C6AF3CB3CAF76433CBBAE6C525AA92D72C15239E1D443FF353AD32F468386D34D058A46BE517A7D02282BD080C29C992635F2AE093E7732FE937DE8A3524391BDF4CE9AC747CF2D060CDDEFD4958473B31E028970B
100000,Account Node,1,044B598464FCB5610BAA175D7D509330056ED7728B59D562BF8DABDBB0954B99,0000000000000000034D494E00C33503A9095C6C96705873B1718BBBB58C9BAC0896CFB5AD82C7BCB4C4B46C3B516E55AE9D90AACE555194EBDF218D9D160BFB452CB0316C67E61C5EB58945583AA513307C765D790BE322300A9C76D849977E718A9DD0359F85B598CD6A14C85AD68FD16C4A3FD91543336D16855CBC30D50E34693FE29C9311ED92F6434449F03C0CDA3652B4AD1FD6D0BB862032B10ABDD84D5BF93ABE06B99A8095AF214128784876C2B6918E28F5549C18079D176E1C5B907E8911D5E34B61D96848D3420000000000000000000000000000000000000000000000000000000000000000D90231527E9DFA2594AA2C32614124DB82252DC6AE28FE62A07DF290DCDDCEB1F28511BFD04B6E41185094C259E5C653D17A431E3463F248EA14E6AC2C694975F09CD4F4681CCACF404A76385A0DDC14F6329339007871D0CB9D506CE1A41155A21E4ECE85E76D0602648E40A23AAC332B9B55F1B9EFDE5F3FDEE49AF83EEAE90000000000000000000000000000000000000000000000000000000000000000EE7EE773D18DE28C9AF12B455F8CF7A5756B438DDFC26F1C90875EC436AE233EDFC31D19CBEB1A67B09362FBAD83DF7C5AAA02BB70D9465676367A2AA5FFBEA500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
100000,AccountRoot,2,C33503A9095C6C96705873B1718BBBB58C9BAC0896CFB5AD82C7BCB4C4B46C3B,0000000000000000034D4C4E001100612200000000240000000A2500017D0A2B3BE715402D00000000413A3C36EEB25EDD249CE94474035CB006551913D115FF45B60E6D0A9A241994D8A3F4A6426B303F0DDAE5A0CDA758C849D262400000003E2AEA66770D776565786368616E67652E636F81147469B8AA28E3EB31ECFFE48C1E66D1BB185DB789006E0B1413DC1D5A8076190562F32B015894CA550AC8A335115F2BF8C453759F
100000,RippleState,2,516E55AE9D90AACE555194EBDF218D9D160BFB452CB0316C67E61C5EB5894558,0000000000000000034D4C4E0011007222000300002500011B3D37000000000000000038000000000000000055921FAE2F6F95DC354122EEFFB86BF894A947BC0E8FCB3B21CB705762B7A5CB7962D4871AFD498D00000000000000000000000000005553440000000000000000000000000000000000000000000000000166D4C38D7EA4C68000000000000000000000000000555344000000000012DC0654E3190F66CC994EF9E214503305B979AD67D4871AFD498D0000000000000000000000000000555344000000000054EE3CE2AC4E9F5524BBCCE0C77F7DEF1CFC46C9017AF787B464E572EE0CEA0777D0D7CE494BC83AEF1F9E97301BF686DFD0B213
100000,Account Node,2,3AA513307C765D790BE322300A9C76D849977E718A9DD0359F85B598CD6A14C8,0000000000000000034D494E0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000073D2969770667D2446BC4FA989074087783941E993FD90C0CF7FE56C8A4D2F460000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000E2A0F55E5274A7B0AE3F075195EB4FF71AFC264322BB29DA33689634F9C011CCCF36E048A0087D50AD6ADA2737FBE1F63FBF1545CDB50FE9E87DF3EB3BC71EA00000000000000000000000000000000000000000000000000000000000000000E09DA149D5C6715F25FFC637C97FC79FD43068DFDF70AFB4CDF637C867A7B4430000000000000000000000000000000000000000000000000000000000000000
...
```

Dump the transactions with columns:

Ledger, Transaction Type, Depth, NodeId, NodeValue

```bash
rdb -path=ripple/nodedb/ -command=transactions
99988,AccountSet,1,56D79CE2BF58F38BB780D64104DEBCCDBD1A1F9786366E6BC26C89A3F826BC1D,000186940001869404534E4400991200032200010000240000003568400000000000000A732102E49205C3DB670A6303AEDE2CACCA4EA4E2C5C38BDB93A5CF7AB17EE7ECC4AAEA74483046022100C8D5D26FDD10DE6F2BBC9B19A67F626134C77A3BF0E8AA92D864EA93BC701667022100A1DC34D6223D5DA46809E464057B1B3129C43DB21E224D06296F269E82A9357E81141B7DAD306A3F274EC898DA985ED9C70B174BDD3997201C00000000F8E51100612500018692553D29D1F990164331FF07647D35D5E116C98FDF948A5604EE097858B22AC337F8560A74A12B225256174DEA8A0725A9990EAE24A9C693DC6BB9572C567E4A695F25E6240000003562400000020BC0D930E1E7220002000024000000362D0000000162400000020BC0D92681141B7DAD306A3F274EC898DA985ED9C70B174BDD39E1E1F1031000A3D2E0B06D57D1521F35CA623FCCAE09BF510185732BDBC306839AA6B25C417F
99986,Payment,1,F9D655AA8722F157B0C6A37E73E96FC05C16FB461C6E136DE8D77980EF5898B7,000186920001869204534E4400B81200002200000000240000003461400000001DCD650068400000000000000A732102E49205C3DB670A6303AEDE2CACCA4EA4E2C5C38BDB93A5CF7AB17EE7ECC4AAEA74483046022100D1A72F916DB5894C199AE478EB642C754A91DDEC7060E6BEAB1E833C5212D442022100DEFAC81FF3F80513CB78F4FC185B468A46884B0181F4BE6E79DB446EA3D3F40781141B7DAD306A3F274EC898DA985ED9C70B174BDD398314C187B1520AEC03B66FD7F1F19D25BF1297C7F4CEC122201C00000000F8E5110061250001868A55831A1EECC9FC06F8B74FB3FA58CBB3CA1A167BAA282D35076B7499954D163E6A560A74A12B225256174DEA8A0725A9990EAE24A9C693DC6BB9572C567E4A695F25E624000000346240000002298E3E3AE1E7220002000024000000352D0000000162400000020BC0D93081141B7DAD306A3F274EC898DA985ED9C70B174BDD39E1E1E311006156D4232D196E9B8DB249D40F1C213CAA72801A5D06D8DDFB1358733AB04DBB5627E8240000000162400000001DCD65008114C187B1520AEC03B66FD7F1F19D25BF1297C7F4CEE1E1F10310003D29D1F990164331FF07647D35D5E116C98FDF948A5604EE097858B22AC337F8
99984,Payment,1,683F23341A3CDE836CB5006D898474A7317D541419E0E6BABE383F8380875961,000186900001869004534E4400B712000022000000002400000027614000000005F5E10068400000000000000A732102CE86AFCD4AB0A4CD25DAF376285ABA3A9CF564C08E69B768E2615B09B886D82B74473045022100E7E5122BD8A0CFC07BFAF86ED98E2FF2D2548AD3DF78AFE25524B9ECFA888955022010963F9FC7E07AF698255E23033B61E2A1DF2DE910A4A1DE0760F0C8621D3B7181145B1CFDD87C0A87B3B1FC0D0D590A787D30765CC18314C187B1520AEC03B66FD7F1F19D25BF1297C7F4CE97201C00000000F8E5110061250001862455EE6B31FA61B24E70812DC9AC76033AF6E13A9A6DEE8A6085FB74D5B35BE690FC56DA8DFF15BA485D874123375FF4E201D20A7D9219A85148B897A4C3D6C63D4C3EE624000000276240000000123B4E8CE1E7220000000024000000282D000000026240000000123B4E8281145B1CFDD87C0A87B3B1FC0D0D590A787D30765CC1E1E1F103107D75F700EF598460E92C48C5B75452162ECB8C60613E94EE80C4FE29CE3BA1E403
99978,Payment,1,183DF0B4F38C1934AF37A6613A0B854106CB3FE74305A1B1337EDCBE03466116,0001868A0001868A04534E4400C11F1200002200000000240000003361D40E35FA931A000000000000000000000000000042544300000000001B7DAD306A3F274EC898DA985ED9C70B174BDD3968400000000000000A732102E49205C3DB670A6303AEDE2CACCA4EA4E2C5C38BDB93A5CF7AB17EE7ECC4AAEA74483046022100EF4882BCC67B00A9AA0E36B6BD8B8BA6F211AAB5DCD3E4AFC9C38FEAF070E8E0022100AC627EB80F36AA604A444A068278ED29470C90045C22408430DD9BD78599D00781141B7DAD306A3F274EC898DA985ED9C70B174BDD398314C187B1520AEC03B66FD7F1F19D25BF1297C7F4CE97201C00000000F8E5110061250001868955D76136E44082018D4EC98CA4D7F415CB9CE04FE5A8140606A6DC04165825AB5C560A74A12B225256174DEA8A0725A9990EAE24A9C693DC6BB9572C567E4A695F25E624000000336240000002298E3E44E1E7220002000024000000342D000000016240000002298E3E3A81141B7DAD306A3F274EC898DA985ED9C70B174BDD39E1E1F103107C831A1EECC9FC06F8B74FB3FA58CBB3CA1A167BAA282D35076B7499954D163E6A
99977,AccountSet,1,D49F066688E38347CFCC66109B42102A3D9CF8FA377561F28E064250CFB0C883,000186890001868904534E4400971200032200010000240000003268400000000000000A732102E49205C3DB670A6303AEDE2CACCA4EA4E2C5C38BDB93A5CF7AB17EE7ECC4AAEA7446304402200E4748EBAB12C53FFC12EB6775B9DFCD42A49748B356FE6D2D74E05DBEC9B22C0220130C89C0DBB2270CF4631FBE3DC86ED8DE6F3D47C7B3D8BE39C0F9258492572981141B7DAD306A3F274EC898DA985ED9C70B174BDD3997201C00000000F8E5110061250001867855D6C7DBBC0B437BB0CED2F4EE1E1503969E5198577198068C4B30F1C1220BD68B560A74A12B225256174DEA8A0725A9990EAE24A9C693DC6BB9572C567E4A695F25E624000000326240000002298E3E4EE1E7220002000024000000332D000000016240000002298E3E4481141B7DAD306A3F274EC898DA985ED9C70B174BDD39E1E1F1031000D76136E44082018D4EC98CA4D7F415CB9CE04FE5A8140606A6DC04165825AB5C
...
```

Dump the account states (leaf nodes only):

Ledger, Ledger Entry Type, Depth, NodeId, NodeValue

```bash
rdb -path=ripple/nodedb/ -command=accounts 
100000,AccountRoot,2,C33503A9095C6C96705873B1718BBBB58C9BAC0896CFB5AD82C7BCB4C4B46C3B,0000000000000000034D4C4E001100612200000000240000000A2500017D0A2B3BE715402D00000000413A3C36EEB25EDD249CE94474035CB006551913D115FF45B60E6D0A9A241994D8A3F4A6426B303F0DDAE5A0CDA758C849D262400000003E2AEA66770D776565786368616E67652E636F81147469B8AA28E3EB31ECFFE48C1E66D1BB185DB789006E0B1413DC1D5A8076190562F32B015894CA550AC8A335115F2BF8C453759F
100000,RippleState,2,516E55AE9D90AACE555194EBDF218D9D160BFB452CB0316C67E61C5EB5894558,0000000000000000034D4C4E0011007222000300002500011B3D37000000000000000038000000000000000055921FAE2F6F95DC354122EEFFB86BF894A947BC0E8FCB3B21CB705762B7A5CB7962D4871AFD498D00000000000000000000000000005553440000000000000000000000000000000000000000000000000166D4C38D7EA4C68000000000000000000000000000555344000000000012DC0654E3190F66CC994EF9E214503305B979AD67D4871AFD498D0000000000000000000000000000555344000000000054EE3CE2AC4E9F5524BBCCE0C77F7DEF1CFC46C9017AF787B464E572EE0CEA0777D0D7CE494BC83AEF1F9E97301BF686DFD0B213
100000,RippleState,3,73D2969770667D2446BC4FA989074087783941E993FD90C0CF7FE56C8A4D2F46,0000000000000000034D4C4E0011007222000200002500016D5137000000000000000038000000000000000055BB17B917AA0BB76EAA9CE4B4C54AFAB1DB5BC5B2204F0CC13F592D25422B30396280000000000000000000000000000000000000004C5443000000000000000000000000000000000000000000000000016680000000000000000000000000000000000000004C5443000000000019820818AA13DF56AB2B4195F9131520CA80461267D5438D7EA4C680000000000000000000000000004C54430000000000F9D8D8EB960BFE54E0DA94854C5DF887995EC485023EA56AC490102DE36C702C3194AE15CBA669399FC113C551851FB15CCA51DC
100000,DirectoryNode,3,E2A0F55E5274A7B0AE3F075195EB4FF71AFC264322BB29DA33689634F9C011CC,0000000000000000034D4C4E001100642200000000364D038D7EA4C680005802BA197D6509F149B9254AE6F09C6DF53A742C36624DA63F4D038D7EA4C6800001110000000000000000000000004254430000000000021165039514E9152E8BC944EE71ABE9D9F0B6F457120311000000000000000000000000000000000000000004110000000000000000000000000000000000000000011320E37E237FF65A4368CE89A80D1671367361B0A74A42F7DB4F20F599E286AB523602BA197D6509F149B9254AE6F09C6DF53A742C36624DA63F4D038D7EA4C68000
100000,AccountRoot,3,CF36E048A0087D50AD6ADA2737FBE1F63FBF1545CDB50FE9E87DF3EB3BC71EA0,0000000000000000034D4C4E00110061220000000024000000032500011F162D0000000055F9648B95E7F36804F30A3F56C749B418AC8918C840901A2A345AFE80218704A96240000000160DC06C8114712B799C79D1EEE3094B59EF9920C7FEB3CE449902CE52E3E46AD340B1C7900F86AFB959AE0C246916E3463905EDD61DE26FFFDD
...
```

Bonus points are awarded for piping into the explain tool:

```bash
go get -u github.com/rubblelabs/ripple/tools/explain
rdb -path=ripple/nodedb/ -command=transactions -dump_format="%[5]X"| explain -
✓ AccountSet  0.00001  ✓   rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH 53       
✓ Payment     0.00001  ✓   rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH => rJeHyfdzw88wbfFQFA7pkyzA81812Fj6Bw 500/XRP                                                      <nil>
✓ Payment     0.00001  ✗   r9JmDFydsvkDK9MtSjxAwYp3NosFWsU62B => rJeHyfdzw88wbfFQFA7pkyzA81812Fj6Bw 100/XRP                                                      <nil>
✓ Payment     0.00001  ✗   rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH => rJeHyfdzw88wbfFQFA7pkyzA81812Fj6Bw 0.04/BTC/rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH                  <nil>
✓ AccountSet  0.00001  ✓   rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH 50       
✓ Payment     0.00001  ½   rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH => rExnTpsgDupeJGW3aK469NhNhPGWGKFDPZ 0.02/BTC/rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH                  <nil>
✓ AccountSet  0.00001  ✓   rsWMoLZhRqTGsmztMRyv5UrE28bTbn8gAH 48  
...
```
