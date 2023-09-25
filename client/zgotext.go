// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package client

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"da": &dictionary{index: daIndex, data: daData},
		"de": &dictionary{index: deIndex, data: deData},
		"en": &dictionary{index: enIndex, data: enData},
		"es": &dictionary{index: esIndex, data: esData},
		"fr": &dictionary{index: frIndex, data: frData},
		"it": &dictionary{index: itIndex, data: itData},
		"nl": &dictionary{index: nlIndex, data: nlData},
		"sl": &dictionary{index: slIndex, data: slData},
		"uk": &dictionary{index: ukIndex, data: ukData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"An error occurred after getting the discovery files for the list of organizations": 1,
	"An error occurred after getting the discovery files for the list of servers":       2,
	"An internal error occurred": 33,
	"Could not retrieve institute access server with URL: '%s' from discovery":                                                       13,
	"Failed to cleanup the VPN connection for the current server":                                                                    26,
	"Failed to get current server for renewing the session":                                                                          28,
	"Failed to get the current server to cleanup the connection":                                                                     24,
	"Failed to set the server with identifier: '%s' as the current":                                                                  22,
	"Failover failed to complete with gateway: '%s' and MTU: '%d'":                                                                   30,
	"Identifier: '%s' for server with type: '%d' is not valid":                                                                       20,
	"Identifier: '%s' for server with type: '%d' is not valid for removal":                                                           23,
	"No secure internet server available to set a location for":                                                                      27,
	"No suitable profiles could be found":                                                                                            9,
	"Profile with ID: '%s' could not be obtained from the server":                                                                    11,
	"Profile with ID: '%s' could not be set":                                                                                         10,
	"The VPN configuration could not be obtained":                                                                                    18,
	"The authorization procedure failed to complete":                                                                                 7,
	"The cleanup process was canceled":                                                                                               25,
	"The client tried to autoconnect to the VPN server: %s, but no secure internet location is found. Please manually connect again": 4,
	"The client tried to autoconnect to the VPN server: %s, but no valid profiles were found. Please manually connect again":         8,
	"The client tried to autoconnect to the VPN server: %s, but you need to authorizate again. Please manually connect again":        6,
	"The current profile could not be found":                                                                                         19,
	"The current server could not be found when getting it for expiry":                                                               3,
	"The custom server with URL: '%s' could not be added":                                                                            17,
	"The identifier that was passed to the library is incorrect":                                                                     12,
	"The institute access server with URL: '%s' could not be added":                                                                  14,
	"The log file with directory: '%s' failed to initialize":                                                                         0,
	"The operation for getting a VPN configuration was canceled":                                                                     21,
	"The renewing process was canceled":                                                                                              29,
	"The secure internet location could not be set":                                                                                  5,
	"The secure internet server with organisation ID: '%s' could not be added":                                                       16,
	"The secure internet server with organisation ID: '%s' could not be retrieved from discovery":                                    15,
	"timeout reached": 31,
	"with cause:":     32,
}

var daIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	// Entry 20 - 3F
	0x00000000, 0x00000000, 0x00000000,
} // Size: 164 bytes

const daData string = ""

var deIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	// Entry 20 - 3F
	0x00000000, 0x00000000, 0x00000000,
} // Size: 164 bytes

const deData string = ""

var enIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x0000003a, 0x0000008c, 0x000000d8,
	0x00000119, 0x0000019b, 0x000001c9, 0x00000244,
	0x00000273, 0x000002ed, 0x00000311, 0x0000033b,
	0x0000037a, 0x000003b5, 0x00000401, 0x00000442,
	0x000004a1, 0x000004ed, 0x00000524, 0x00000550,
	0x00000577, 0x000005b6, 0x000005f1, 0x00000632,
	0x0000067d, 0x000006b8, 0x000006d9, 0x00000715,
	0x0000074f, 0x00000785, 0x000007a7, 0x000007ea,
	// Entry 20 - 3F
	0x000007fa, 0x00000806, 0x00000821,
} // Size: 164 bytes

const enData string = "" + // Size: 2081 bytes
	"\x02The log file with directory: '%[1]s' failed to initialize\x02An erro" +
	"r occurred after getting the discovery files for the list of organizatio" +
	"ns\x02An error occurred after getting the discovery files for the list o" +
	"f servers\x02The current server could not be found when getting it for e" +
	"xpiry\x02The client tried to autoconnect to the VPN server: %[1]s, but n" +
	"o secure internet location is found. Please manually connect again\x02Th" +
	"e secure internet location could not be set\x02The client tried to autoc" +
	"onnect to the VPN server: %[1]s, but you need to authorizate again. Plea" +
	"se manually connect again\x02The authorization procedure failed to compl" +
	"ete\x02The client tried to autoconnect to the VPN server: %[1]s, but no " +
	"valid profiles were found. Please manually connect again\x02No suitable " +
	"profiles could be found\x02Profile with ID: '%[1]s' could not be set\x02" +
	"Profile with ID: '%[1]s' could not be obtained from the server\x02The id" +
	"entifier that was passed to the library is incorrect\x02Could not retrie" +
	"ve institute access server with URL: '%[1]s' from discovery\x02The insti" +
	"tute access server with URL: '%[1]s' could not be added\x02The secure in" +
	"ternet server with organisation ID: '%[1]s' could not be retrieved from " +
	"discovery\x02The secure internet server with organisation ID: '%[1]s' co" +
	"uld not be added\x02The custom server with URL: '%[1]s' could not be add" +
	"ed\x02The VPN configuration could not be obtained\x02The current profile" +
	" could not be found\x02Identifier: '%[1]s' for server with type: '%[2]d'" +
	" is not valid\x02The operation for getting a VPN configuration was cance" +
	"led\x02Failed to set the server with identifier: '%[1]s' as the current" +
	"\x02Identifier: '%[1]s' for server with type: '%[2]d' is not valid for r" +
	"emoval\x02Failed to get the current server to cleanup the connection\x02" +
	"The cleanup process was canceled\x02Failed to cleanup the VPN connection" +
	" for the current server\x02No secure internet server available to set a " +
	"location for\x02Failed to get current server for renewing the session" +
	"\x02The renewing process was canceled\x02Failover failed to complete wit" +
	"h gateway: '%[1]s' and MTU: '%[2]d'\x02timeout reached\x02with cause:" +
	"\x02An internal error occurred"

var esIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x0000004a, 0x000000ab, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	0x00000104, 0x00000104, 0x00000104, 0x00000104,
	// Entry 20 - 3F
	0x00000104, 0x00000104, 0x00000104,
} // Size: 164 bytes

const esData string = "" + // Size: 260 bytes
	"\x02El archivo de registro con el directorio: '%[1]s' no se puede inicia" +
	"lizar\x02Se ha producido un error al obtener los archivos de detección d" +
	"e la lista de las organizaciones\x02Se ha producido un error al obtener " +
	"los archivos de detección de la lista de servidores"

var frIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x0000004f, 0x000000ba, 0x00000120,
	0x00000120, 0x000001de, 0x0000021d, 0x0000021d,
	0x0000021d, 0x000002c5, 0x000002ed, 0x00000324,
	0x00000365, 0x00000399, 0x00000399, 0x00000399,
	0x00000399, 0x000003f4, 0x0000043b, 0x00000469,
	0x00000490, 0x000004db, 0x000004db, 0x000004db,
	0x000004db, 0x000004db, 0x000004db, 0x000004db,
	0x000004db, 0x000004db, 0x000004db, 0x000004db,
	// Entry 20 - 3F
	0x000004db, 0x000004db, 0x000004db,
} // Size: 164 bytes

const frData string = "" + // Size: 1243 bytes
	"\x02Le fichier de registre du répertoire\u202f: '%[1]s' n'a pas pu être " +
	"initialisé\x02Une erreur est survenue pendant la récupération des fichie" +
	"rs de détection de la liste des organisations\x02Une erreur est survenue" +
	" pendant la récupération des fichiers de détection de la liste des serve" +
	"urs\x02Le client a essayé de se connecter automatiquement au serveur VPN" +
	" : %[1]s, mais aucune localisation internet sécurisée n'a été trouvée. V" +
	"euillez vous connecter manuellement de nouveau\x02La localisation intern" +
	"et sécurisée n'a pas pu être définie\x02Le client a essayé de se connect" +
	"er automatiquement au serveur VPN : %[1]s, mais aucun profil valide n'a " +
	"été trouvé. Veuillez vous connecter manuellement de nouveau\x02Aucun pr" +
	"ofil adéquat n'a été trouvé\x02Le profil avec l'ID : '%[1]s' n'a pas pu " +
	"être défini\x02Le profil avec l'ID : '%[1]s' n'a pas pu être obtenu du " +
	"serveur\x02L'identifiant envoyé à la librairie est incorrect\x02Le serve" +
	"ur internet sécurisé avec l'ID d'organisation : '%[1]s' n'a pas pu être " +
	"ajouté\x02Le serveur personnalisé avec l'URL : '%[1]s' n'a pas pu être a" +
	"jouté\x02La configuration VPN n'a pas pu être obtenue\x02Le profil actue" +
	"l n'a pas été trouvé\x02L'identifiant : '%[1]s' du serveur avec le type " +
	": '%[2]d' n'est pas valide"

var itIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	// Entry 20 - 3F
	0x00000000, 0x00000000, 0x00000000,
} // Size: 164 bytes

const itData string = ""

var nlIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x0000003c, 0x00000084, 0x000000c7,
	0x000000f6, 0x00000180, 0x000001be, 0x0000023a,
	0x0000026b, 0x000002f5, 0x00000329, 0x00000368,
	0x000003b4, 0x000003f2, 0x00000446, 0x0000048d,
	0x000004e7, 0x00000534, 0x0000056a, 0x00000598,
	0x000005c4, 0x00000608, 0x00000646, 0x00000693,
	0x000006ec, 0x00000737, 0x00000737, 0x00000763,
	0x000007b4, 0x000007fb, 0x000007fb, 0x0000084c,
	// Entry 20 - 3F
	0x0000089d, 0x000008aa, 0x000008c9,
} // Size: 164 bytes

const nlData string = "" + // Size: 2249 bytes
	"\x02Het log bestand met pad: '%[1]s' kan niet aangemaakt worden\x02Er is" +
	" een fout opgetreden met het ophalen van de lijst van organisaties\x02Er" +
	" is een fout opgetreden met het ophalen van de lijst van servers\x02De h" +
	"uidige VPN server kon niet worden gevonden\x02De client wilde automatisc" +
	"h verbinden met VPN server: %[1]s, maar de huidige locatie is niet gevon" +
	"den. U moet opnieuw handmatig verbinden\x02De locatie voor de secure int" +
	"ernet server kon niet opgeslagen\x02De client wilde automatisch verbinde" +
	"n met VPN server: %[1]s, maar authorizatie is nodig. U moet handmatig op" +
	"nieuw verbinden\x02Het authorizatie proces kon niet voltooid worden\x02D" +
	"e client wilde automatisch verbinden met VPN server: %[1]s, maar er was " +
	"geen geldig profiel gevinden. U moet handmatig opnieuw verbinden\x02Er z" +
	"ijn geen profielen gevonden om mee te verbinden\x02Het profiel met ident" +
	"iteit: '%[1]s' kon niet opgeslagen worden\x02Het profiel met identiteit:" +
	" '%[1]s' kon niet opgehaald worden van de server\x02De identiteit die aa" +
	"n de library werd gegeven is niet correct\x02De institute access server " +
	"met URL: '%[1]s' kan niet opgehaald worden van discovery\x02De institute" +
	" access server met URL: '%[1]s' kan niet toegevoegd worden\x02De secure " +
	"internet server met identiteit: '%[1]s' kan niet opgehaald worden van di" +
	"scovery\x02De secure internet server met identiteit: '%[1]s' kan niet to" +
	"egevoegd worden\x02De server met URL: '%[1]s' kan niet toegevoegd worden" +
	"\x02De VPN configuratie kan niet opgehaald worden\x02Het huidig profiel " +
	"kan niet gevonden worden\x02De identiteit: '%[1]s' voor server met type:" +
	" '%[2]d' is niet geldig\x02De procedure om een VPN configuratie op te ha" +
	"len is verbroken\x02De server met identiteit: '%[1]s' kan niet als de hu" +
	"dige server gezet worden\x02Identiteit: '%[1]s' voor server met type: '%" +
	"[2]d' is niet geldig om verwijderd te worden\x02De huidige server kan ni" +
	"et opgehaald worden om de connectie te verwijderen\x02De VPN connectie i" +
	"s niet volledig opgeruimd\x02Er is geen \x22secure internet\x22 server b" +
	"eschikbaar om de locatie voor in te stellen\x02De huidige server kan nie" +
	"t opgehaald worden om de sessie te hernieuwen\x02Het 'failover' proces k" +
	"an niet voltooid worden. Gateway: '%[1]s' en MTU: '%[2]d'\x02Er is een t" +
	"ime-out opgetreden in de verbinding. Controleer uw internetverbinding" +
	"\x02met oorzaak:\x02Een interne fout is opgetreden"

var slIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x0000003c, 0x0000007c, 0x000000bc,
	0x000000f8, 0x0000019b, 0x000001cb, 0x00000264,
	0x00000285, 0x00000325, 0x0000034e, 0x00000381,
	0x000003c1, 0x000003e2, 0x0000043a, 0x00000484,
	0x000004dc, 0x0000052a, 0x00000566, 0x00000588,
	0x000005b0, 0x000005e2, 0x0000061b, 0x00000649,
	0x00000688, 0x000006cb, 0x000006cb, 0x00000704,
	0x00000746, 0x00000788, 0x00000788, 0x000007bf,
	// Entry 20 - 3F
	0x000007cf, 0x000007d9, 0x000007f7,
} // Size: 164 bytes

const slData string = "" + // Size: 2039 bytes
	"\x02Napaka pri vzpostavitvi datoteke dnevnika v imeniku '%[1]s'\x02Pri n" +
	"alaganju datotek kataloga organizacij je prišlo do napake\x02Pri nalagan" +
	"ju datotek kataloga strežnikov je prišlo do napake\x02Ugotavljanje prete" +
	"ka ne more določiti izbranega strežnika\x02Odjemalec je poskusil samodej" +
	"no vzpostaviti povezavo s strežnikom VPN %[1]s, vendar ni našel nobene l" +
	"okacije za varni splet. Ponovno vzpostavite povezavo ročno\x02Napaka pri" +
	" nastavljanju lokacije za varni splet\x02Odjemalec je poskusil samodejno" +
	" vzpostaviti povezavo s strežnikom VPN %[1]s, vendar ga morate ponovno a" +
	"vtorizirati. Ponovno vzpostavite povezavo ročno\x02Napaka pri postopku a" +
	"vtorizacije\x02Odjemalec je poskusil samodejno vzpostaviti povezavo s st" +
	"režnikom VPN %[1]s, vendar ni našel nobenega veljavnega profila. Ponovno" +
	" vzpostavite povezavo ročno\x02Ustreznih profilov ni bilo mogoče najti" +
	"\x02Profila z ID-jem '%[1]s' ni bilo mogoče nastaviti\x02Profila z ID-je" +
	"m '%[1]s' ni bilo mogoče naložiti s strežnika\x02ID poslan knjižnici je " +
	"napačen\x02Strežnika z naslovom URL '%[1]s' za dostop do ustanove ni bil" +
	"o možno najti v katalogu\x02Strežnika z naslovom '%[1]s' za dostop do us" +
	"tanove ni bilo možno dodati\x02Strežnika za varni splet organizacije z I" +
	"D-jem '%[1]s' ni bilo možno najti v katalogu\x02Strežnika za varni splet" +
	" organizacije z ID-jem '%[1]s' ni bilo možno dodati\x02Svojega strežnika" +
	" z naslovom '%[1]s' ni bilo možno dodati\x02Napaka pri prenosu nastavite" +
	"v VPN\x02Izbranega profila ni bilo mogoče najti\x02ID '%[1]s' za strežni" +
	"k vrste '%[2]d' ni veljaven\x02Operacija za nalaganje nastavitev VPN je " +
	"bila preklicana\x02Napaka pri izbiri strežnika z ID-jem '%[1]s'\x02ID '%" +
	"[1]s' strežnika vrste '%[2]d' ni veljaven za odstranitev\x02Napaka pri d" +
	"oločanju izbranega strežnika za čiščenje povezave\x02Napaka pri čiščenju" +
	" povezave VPN za izbrani strežnik\x02Za izbiro lokacije ni na voljo nobe" +
	"nega strežnika za varni splet\x02Napaka pri ugotavljanju izbranega strež" +
	"nika za podaljšanje seje\x02Preklop ni uspel s prehodom '%[1]s' in MTU-j" +
	"em '%[2]d'\x02čas je potekel\x02; razlog:\x02Prišlo je do notranje napak" +
	"e"

var ukIndex = []uint32{ // 35 elements
	// Entry 0 - 1F
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	0x00000000, 0x00000000, 0x00000000, 0x00000000,
	// Entry 20 - 3F
	0x00000000, 0x00000000, 0x00000000,
} // Size: 164 bytes

const ukData string = ""

// Total table size 9348 bytes (9KiB); checksum: 31AC92B
