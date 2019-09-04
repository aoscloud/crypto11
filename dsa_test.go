// Copyright 2016, 2017 Thales e-Security, Inc
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package crypto11

import (
	"crypto"
	"crypto/dsa"
	"crypto/rand"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

// Use pre-cooked groups, making new ones is too slow and doesn't test
// the thing we care about right now
var dsaSizes = map[dsa.ParameterSizes]*dsa.Parameters{
	dsa.L1024N160: dsaParameters("162901807344330501683198681796158038892896087086101563643117520955342956222658660947046551471908694917426639943346403004016664903097104319514472496004393646085043684550200397339750548531631673569299660859981509826530366281432841569092450374450801368309873677511310151482654326760484535441347619035325248971027",
		"1152435081592794070452857956121091711775772698753",
		"115452448594884027370996574512176932912212827128999848033626686707424652518655228807676153362807773892059885310384996519508710450365985083911488963280781008704774827093959271338258878644127595592794116557826253235239116215434665533711823763161327640232154970352210500147125324244313495378864886440434350445047"),
	dsa.L2048N224: dsaParameters("26897309010050677376706120375514893443234451462311845264107450103587616694552775798203063987684778976094649425235043381904342317828282144890959162115530648385936228003908801487852374144335466990519500032005785121648159893054617696809139702264373725369373070626941633457622971603832559810490847598977955584343798411835596790144008113895279879074659082117447812649062019962590823347376506566929271773368811020476842094620754802895105872616973988885001858825172067780326607348935255985451124288300256500395655426867070867964116873929476668281059072072214310998391443703378306986615500676546104049539495719703519940117951",
		"21279492449681152694989309770053968015769408476251418501127004412991",
		"806074083423550090996520226918161019531020530678027294640523175956127336741123734330141199877500196483993261275159586661624671147930168776905666558638413092940255649594308870768538528038043504442900678005068039885085228576696098584188275093427119292669713986799663707376358075013805220640084303178400319317901835943666702067023936970119214055180476824871264924202941665618311395109943406576227418395899047272659412653021884375802245946147395739665170684226546861490208766142484366132284564686737604148267296225461433426617875614890925967933759512314096027440880608838689527799562179981619366217116863163848643347517"),
	dsa.L2048N256: dsaParameters("17112618140558405746569134228451035632946851190522147937285442098259988182250197615427085779160410421439180425324757684372143777133652794766735361927034003821174405167502004454639027705200074749864396322138822913233787614500370035203936072693687262723523490964341964969759901927635271723521044671768174187470434319243240992425833709397541419254550881926245866992880605910541826307493251693926611230487275444183879966706699255763528432543168158080992127766526692139601189110012363846432577767973841827921723244473497325201586745757496510772406382623383426532181890987571204825714955203879088898058005561306135204877009",
		"91121968019586848097218691760416459853586419261325465835456057472782359091757",
		"11440584645174460641968590902479650665148631776226218292071277461507925953586967572345732823684968840767915723600616468676971247641561301128304315573382612091463182732020028761013121095491181166498021161181825176128544846296532345606186214409027298873946181476041031383608259830471735783109449938372605549135088084178714254363472647038149745329509161045818621699799741032979200140196670711384411636885603577783522623156072149396954757319730766808882432057418486340662640743658615448774457850848807133015183403659298100852797177921099769325832019636521852355306536644565775804923805884803806962749273871634832496124563"),
	dsa.L3072N256: dsaParameters("3292223093076756007575154742534996337119773908047183722225998562022708727169602325355570673703977126031767898507744907641046577256340035141288034004031488512406489120451045534249225732833971818604693952828746803787139641400035401666093696260981834564456248118055149793765296502054830103686831403950712801607640861048795967036349397959545638669305551590712409193326920723566244393566858414736680142461926772620503656504255924167986033103929106720448272225556256931951085963366484599500380306971899408684747849004037207957694189418415781274826562470103438926910486414615463061843381314235226091370587221926283099048619624319791834455498858310112354321475086722140598386575332287897366827250193444772749165320615259571588848174295398383162626179797293184955390849951880022388527593205033454738502281212004414696356545085156516647817931381822225669880911822774112408110220283747435924461722168415951571106537524166169188892805803",
		"108867360151102455324355349853282813937921700957084486520964693771814677113683",
		"1557152939324444301267588704765106652107942442192133066719901270619999543077285434937381472603966361780723031759112891413917740586179484242767473407407437414309703817996415200985589914186423758047030899237374015161656689079147160316057703866877939220177527975058770672610222076376028225539948249840143348317913168294956877608194868405476013052838117032232844130922725529104417156431908522534312701383683651539112438656112551665554183084212392173005974168567420494917967524105375970433753993335902159248975586830245477125927759707545151441829187496420494441219127113977851955484097152895793192460180527358969782742628457686090239586702325662876026163849101009626548110941786173048851538499449581500118606509607658612519792319018156890523385946253475322407270910172103139620247134330065127064271570620806081501157521483092036300763874199611913859815207922006817632837708771028001952742328571406294474170839325153987075519613561"),
}

func dsaParameters(p, q, g string) *dsa.Parameters {
	var params = dsa.Parameters{
		P: new(big.Int),
		Q: new(big.Int),
		G: new(big.Int),
	}
	params.P.SetString(p, 10)
	params.Q.SetString(q, 10)
	params.G.SetString(g, 10)
	return &params
}

// Make dsa.PrivateKey look like a Crypto.Signer

type DSAPrivateKey dsa.PrivateKey

func (signer *DSAPrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	key := (*dsa.PrivateKey)(signer)
	var sig dsaSignature
	var err error
	if sig.R, sig.S, err = dsa.Sign(rand, key, digest); err != nil {
		return nil, err
	}
	return sig.marshalDER()
}

func (signer *DSAPrivateKey) Public() crypto.PublicKey {
	return crypto.PublicKey(&signer.PublicKey)
}

func TestNativeDSA(t *testing.T) {
	var err error
	for psize, params := range dsaSizes {
		var key = &dsa.PrivateKey{
			PublicKey: dsa.PublicKey{
				Parameters: *params,
				Y:          nil,
			},
			X: nil,
		}
		if err = dsa.GenerateKey(key, rand.Reader); err != nil {
			t.Errorf("crypto.dsa.GenerateKey: %v", err)
			return
		}
		testDsaSigning(t, (*DSAPrivateKey)(key), psize, "native")
	}
}

func TestHardDSA(t *testing.T) {
	ctx, err := ConfigureFromFile("config")
	require.NoError(t, err)

	defer func() {
		err = ctx.Close()
		require.NoError(t, err)
	}()

	for pSize, params := range dsaSizes {

		id := randomBytes()
		label := randomBytes()

		key, err := ctx.GenerateDSAKeyPairWithLabel(id, label, params)
		require.NoError(t, err, "Failed for key size %s", parameterSizeToString(pSize))
		defer func() { _ = key.Delete() }()

		testDsaSigning(t, key, pSize, "hard1")

		key2, err := ctx.FindKeyPair(id, nil)
		require.NoError(t, err)
		testDsaSigning(t, key2.(*pkcs11PrivateKeyDSA), pSize, "hard2")

		key3, err := ctx.FindKeyPair(nil, label)
		require.NoError(t, err)
		testDsaSigning(t, key3.(crypto.Signer), pSize, "hard3")
	}
}

func parameterSizeToString(s dsa.ParameterSizes) string {
	switch s {
	case dsa.L1024N160:
		return "L1024N160"
	case dsa.L2048N224:
		return "L2048N224"
	case dsa.L2048N256:
		return "L2048N256"
	case dsa.L3072N256:
		return "L3072N256"
	default:
		return "unknown"
	}
}

func testDsaSigning(t *testing.T, key crypto.Signer, psize dsa.ParameterSizes, what string) {
	testDsaSigningWithHash(t, key, crypto.SHA1, psize, what)
	testDsaSigningWithHash(t, key, crypto.SHA224, psize, what)
	testDsaSigningWithHash(t, key, crypto.SHA256, psize, what)
	testDsaSigningWithHash(t, key, crypto.SHA384, psize, what)
	testDsaSigningWithHash(t, key, crypto.SHA512, psize, what)
}

func testDsaSigningWithHash(t *testing.T, key crypto.Signer, hashFunction crypto.Hash, psize dsa.ParameterSizes, what string) {

	plaintext := []byte("sign me with DSA")
	h := hashFunction.New()
	_, err := h.Write(plaintext)
	require.NoError(t, err)

	plaintextHash := h.Sum([]byte{}) // weird API
	// crypto.dsa.Sign doesn't truncate the hash!
	qbytes := (dsaSizes[psize].Q.BitLen() + 7) / 8
	plaintextHash = plaintextHash[:qbytes]

	sigDER, err := key.Sign(rand.Reader, plaintextHash, hashFunction)
	require.NoError(t, err)

	var sig dsaSignature
	err = sig.unmarshalDER(sigDER)
	require.NoError(t, err)

	dsaPubkey := key.Public().(crypto.PublicKey).(*dsa.PublicKey)
	if !dsa.Verify(dsaPubkey, plaintextHash, sig.R, sig.S) {
		t.Errorf("DSA %s Verify failed (psize %d hash %v)", what, psize, hashFunction)
	}
}

func TestDsaRequiredArgs(t *testing.T) {
	ctx, err := ConfigureFromFile("config")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, ctx.Close())
	}()

	_, err = ctx.GenerateDSAKeyPair(nil, dsaSizes[dsa.L2048N224])
	require.Error(t, err)

	val := randomBytes()

	_, err = ctx.GenerateDSAKeyPairWithLabel(nil, val, dsaSizes[dsa.L2048N224])
	require.Error(t, err)

	_, err = ctx.GenerateDSAKeyPairWithLabel(val, nil, dsaSizes[dsa.L2048N224])
	require.Error(t, err)
}
