// Copyright (c) 2017 Ian Davis
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

// See https://en.bitcoin.it/wiki/Wallet_import_format for details on the encoding/decoding of private keys

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "  %s [KEY]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "    Encode a key into a word list")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "  %s [WORD] [WORD] ...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "    Decode a word list into a key")
		fmt.Fprintln(os.Stderr)
	}
	flag.Parse()
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Main() error {
	switch flag.NArg() {
	case 0:
		return fmt.Errorf("expected one or more arguments")
	case 1:
		return encode(flag.Arg(0))
	default:
		return decode(flag.Args())
	}
}

func encode(s string) error {
	key, version, err := CheckDecode(s)
	if err != nil {
		return err
	}

	if len(key) != 32 {
		return fmt.Errorf("key is wrong length, got %d, wanted 32 bytes", len(key))
	}

	fmt.Printf("%x version %d\n", key, version)
	even := true
	for _, b := range key {
		if even {
			fmt.Print(words[int(b)].even, " ")
		} else {
			fmt.Print(words[int(b)].odd, " ")
		}
		even = !even
	}
	fmt.Println()
	return nil

}

func decode(wds []string) error {
	if len(wds) != 32 {
		return fmt.Errorf("wrong number of words, got %d, wanted 32", len(wds))
	}

	evenWords := make(map[string]byte, 256)
	oddWords := make(map[string]byte, 256)
	for i := range words {
		evenWords[words[i].even] = byte(i)
		oddWords[words[i].odd] = byte(i)
	}

	key := make([]byte, 32)
	even := true
	for i := 0; i < 32; i++ {
		if even {
			b, ok := evenWords[wds[i]]
			if !ok {
				if _, ok := oddWords[wds[i]]; ok {
					return fmt.Errorf("invalid word: found %s in odd position when expecting an odd word. Is word list truncated?", wds[i])
				}
				return fmt.Errorf("invalid word: %s is not in the word list", wds[i])
			}
			key[i] = b
		} else {
			b, ok := oddWords[wds[i]]
			if !ok {
				if _, ok := evenWords[wds[i]]; ok {
					return fmt.Errorf("invalid word: found %s in even position when expecting an odd word. Is word list truncated?", wds[i])
				}
				return fmt.Errorf("invalid word: %s is not in the word list", wds[i])
			}
			key[i] = b
		}
		even = !even
	}
	fmt.Printf("Hex: %x\n", key)
	encoded := CheckEncode(key, 128)
	fmt.Printf("Base58Check encoded: %s\n", encoded)
	return nil
}

// words are the PGP Word List which are two lists of 256 words for even and
// odd input bytes. The even words have two syllables whereas the odd words
// have three. See https://en.wikipedia.org/wiki/PGP_word_list and https://we
// b.archive.org/web/20080303003708/http://web.mit.edu:80/network/pgpfone/manu
// al/index.html
var words = []struct {
	even, odd string
}{
	{"aardvark", "adroitness"},
	{"absurd", "adviser"},
	{"accrue", "aftermath"},
	{"acme", "aggregate"},
	{"adrift", "alkali"},
	{"adult", "almighty"},
	{"afflict", "amulet"},
	{"ahead", "amusement"},
	{"aimless", "antenna"},
	{"Algol", "applicant"},
	{"allow", "Apollo"},
	{"alone", "armistice"},
	{"ammo", "article"},
	{"ancient", "asteroid"},
	{"apple", "Atlantic"},
	{"artist", "atmosphere"},
	{"assume", "autopsy"},
	{"Athens", "Babylon"},
	{"atlas", "backwater"},
	{"Aztec", "barbecue"},
	{"baboon", "belowground"},
	{"backfield", "bifocals"},
	{"backward", "bodyguard"},
	{"banjo", "bookseller"},
	{"beaming", "borderline"},
	{"bedlamp", "bottomless"},
	{"beehive", "Bradbury"},
	{"beeswax", "bravado"},
	{"befriend", "Brazilian"},
	{"Belfast", "breakaway"},
	{"berserk", "Burlington"},
	{"billiard", "businessman"},
	{"bison", "butterfat"},
	{"blackjack", "Camelot"},
	{"blockade", "candidate"},
	{"blowtorch", "cannonball"},
	{"bluebird", "Capricorn"},
	{"bombast", "caravan"},
	{"bookshelf", "caretaker"},
	{"brackish", "celebrate"},
	{"breadline", "cellulose"},
	{"breakup", "certify"},
	{"brickyard", "chambermaid"},
	{"briefcase", "Cherokee"},
	{"Burbank", "Chicago"},
	{"button", "clergyman"},
	{"buzzard", "coherence"},
	{"cement", "combustion"},
	{"chairlift", "commando"},
	{"chatter", "company"},
	{"checkup", "component"},
	{"chisel", "concurrent"},
	{"choking", "confidence"},
	{"chopper", "conformist"},
	{"Christmas", "congregate"},
	{"clamshell", "consensus"},
	{"classic", "consulting"},
	{"classroom", "corporate"},
	{"cleanup", "corrosion"},
	{"clockwork", "councilman"},
	{"cobra", "crossover"},
	{"commence", "crucifix"},
	{"concert", "cumbersome"},
	{"cowbell", "customer"},
	{"crackdown", "Dakota"},
	{"cranky", "decadence"},
	{"crowfoot", "December"},
	{"crucial", "decimal"},
	{"crumpled", "designing"},
	{"crusade", "detector"},
	{"cubic", "detergent"},
	{"dashboard", "determine"},
	{"deadbolt", "dictator"},
	{"deckhand", "dinosaur"},
	{"dogsled", "direction"},
	{"dragnet", "disable"},
	{"drainage", "disbelief"},
	{"dreadful", "disruptive"},
	{"drifter", "distortion"},
	{"dropper", "document"},
	{"drumbeat", "embezzle"},
	{"drunken", "enchanting"},
	{"Dupont", "enrollment"},
	{"dwelling", "enterprise"},
	{"eating", "equation"},
	{"edict", "equipment"},
	{"egghead", "escapade"},
	{"eightball", "Eskimo"},
	{"endorse", "everyday"},
	{"endow", "examine"},
	{"enlist", "existence"},
	{"erase", "exodus"},
	{"escape", "fascinate"},
	{"exceed", "filament"},
	{"eyeglass", "finicky"},
	{"eyetooth", "forever"},
	{"facial", "fortitude"},
	{"fallout", "frequency"},
	{"flagpole", "gadgetry"},
	{"flatfoot", "Galveston"},
	{"flytrap", "getaway"},
	{"fracture", "glossary"},
	{"framework", "gossamer"},
	{"freedom", "graduate"},
	{"frighten", "gravity"},
	{"gazelle", "guitarist"},
	{"Geiger", "hamburger"},
	{"glitter", "Hamilton"},
	{"glucose", "handiwork"},
	{"goggles", "hazardous"},
	{"goldfish", "headwaters"},
	{"gremlin", "hemisphere"},
	{"guidance", "hesitate"},
	{"hamlet", "hideaway"},
	{"highchair", "holiness"},
	{"hockey", "hurricane"},
	{"indoors", "hydraulic"},
	{"indulge", "impartial"},
	{"inverse", "impetus"},
	{"involve", "inception"},
	{"island", "indigo"},
	{"jawbone", "inertia"},
	{"keyboard", "infancy"},
	{"kickoff", "inferno"},
	{"kiwi", "informant"},
	{"klaxon", "insincere"},
	{"locale", "insurgent"},
	{"lockup", "integrate"},
	{"merit", "intention"},
	{"minnow", "inventive"},
	{"miser", "Istanbul"},
	{"Mohawk", "Jamaica"},
	{"mural", "Jupiter"},
	{"music", "leprosy"},
	{"necklace", "letterhead"},
	{"Neptune", "liberty"},
	{"newborn", "maritime"},
	{"nightbird", "matchmaker"},
	{"Oakland", "maverick"},
	{"obtuse", "Medusa"},
	{"offload", "megaton"},
	{"optic", "microscope"},
	{"orca", "microwave"},
	{"payday", "midsummer"},
	{"peachy", "millionaire"},
	{"pheasant", "miracle"},
	{"physique", "misnomer"},
	{"playhouse", "molasses"},
	{"Pluto", "molecule"},
	{"preclude", "Montana"},
	{"prefer", "monument"},
	{"preshrunk", "mosquito"},
	{"printer", "narrative"},
	{"prowler", "nebula"},
	{"pupil", "newsletter"},
	{"puppy", "Norwegian"},
	{"python", "October"},
	{"quadrant", "Ohio"},
	{"quiver", "onlooker"},
	{"quota", "opulent"},
	{"ragtime", "Orlando"},
	{"ratchet", "outfielder"},
	{"rebirth", "Pacific"},
	{"reform", "pandemic"},
	{"regain", "Pandora"},
	{"reindeer", "paperweight"},
	{"rematch", "paragon"},
	{"repay", "paragraph"},
	{"retouch", "paramount"},
	{"revenge", "passenger"},
	{"reward", "pedigree"},
	{"rhythm", "Pegasus"},
	{"ribcage", "penetrate"},
	{"ringbolt", "perceptive"},
	{"robust", "performance"},
	{"rocker", "pharmacy"},
	{"ruffled", "phonetic"},
	{"sailboat", "photograph"},
	{"sawdust", "pioneer"},
	{"scallion", "pocketful"},
	{"scenic", "politeness"},
	{"scorecard", "positive"},
	{"Scotland", "potato"},
	{"seabird", "processor"},
	{"select", "provincial"},
	{"sentence", "proximate"},
	{"shadow", "puberty"},
	{"shamrock", "publisher"},
	{"showgirl", "pyramid"},
	{"skullcap", "quantity"},
	{"skydive", "racketeer"},
	{"slingshot", "rebellion"},
	{"slowdown", "recipe"},
	{"snapline", "recover"},
	{"snapshot", "repellent"},
	{"snowcap", "replica"},
	{"snowslide", "reproduce"},
	{"solo", "resistor"},
	{"southward", "responsive"},
	{"soybean", "retraction"},
	{"spaniel", "retrieval"},
	{"spearhead", "retrospect"},
	{"spellbind", "revenue"},
	{"spheroid", "revival"},
	{"spigot", "revolver"},
	{"spindle", "sandalwood"},
	{"spyglass", "sardonic"},
	{"stagehand", "Saturday"},
	{"stagnate", "savagery"},
	{"stairway", "scavenger"},
	{"standard", "sensation"},
	{"stapler", "sociable"},
	{"steamship", "souvenir"},
	{"sterling", "specialist"},
	{"stockman", "speculate"},
	{"stopwatch", "stethoscope"},
	{"stormy", "stupendous"},
	{"sugar", "supportive"},
	{"surmount", "surrender"},
	{"suspense", "suspicious"},
	{"sweatband", "sympathy"},
	{"swelter", "tambourine"},
	{"tactics", "telephone"},
	{"talon", "therapist"},
	{"tapeworm", "tobacco"},
	{"tempest", "tolerance"},
	{"tiger", "tomorrow"},
	{"tissue", "torpedo"},
	{"tonic", "tradition"},
	{"topmost", "travesty"},
	{"tracker", "trombonist"},
	{"transit", "truncated"},
	{"trauma", "typewriter"},
	{"treadmill", "ultimate"},
	{"Trojan", "undaunted"},
	{"trouble", "underfoot"},
	{"tumor", "unicorn"},
	{"tunnel", "unify"},
	{"tycoon", "universe"},
	{"uncut", "unravel"},
	{"unearth", "upcoming"},
	{"unwind", "vacancy"},
	{"uproot", "vagabond"},
	{"upset", "vertigo"},
	{"upshot", "Virginia"},
	{"vapor", "visitor"},
	{"village", "vocalist"},
	{"virus", "voyager"},
	{"Vulcan", "warranty"},
	{"waffle", "Waterloo"},
	{"wallet", "whimsical"},
	{"watchword", "Wichita"},
	{"wayside", "Wilmington"},
	{"willow", "Wyoming"},
	{"woodlark", "yesteryear"},
	{"Zulu", "Yucatan"},
}
