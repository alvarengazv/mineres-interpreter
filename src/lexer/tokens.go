package lexer

type TabelaPalavras int

const (
	trem_di_numeru      TabelaPalavras = iota // 0
	trem_cum_virgula                          // 1
	trem_discrita                             // 2
	trem_discolhe                             // 3
	trosso                                    // 4
	uai_se                                    // 5
	uai_senao                                 // 6
	roda_esse_trem                            // 7
	enquanto_tiver_trem                       // 8
	dependenu                                 // 9
	du_casu                                   // 10
	ta_bao                                    // 11
	para_o_trem                               // 12
	toca_o_trem                               // 13
	bora_cumpade                              // 14
	abre_parentese                            // 15
	fecha_parentese                           // 16
	eh                                        // 17
	num_eh                                    // 18
	simbora                                   // 19
	cabo                                      // 20
	abre_chave                                // 21
	fecha_chave                               // 22
	abre_aspas                                // 23z
	fecha_aspas                               // 24
	virgula                                   // 25
	uai                                       // 26
	menor_que                                 // 27
	maior_que                                 // 28
	menor_igual_que                           // 29
	maior_igual_que                           // 30
	fica_assim_entao                          // 31
	neh_nada                                  // 32
	mema_coisa                                // 33
	quarque_um                                // 34
	vam_marca                                 // 35
	tamem                                     // 36
	um_o_oto                                  // 37
	soma                                      // 38
	subtracao                                 // 39
	veiz                                      // 40
	sob                                       // 41
	modulo                                    // 42
	divisao_inteira                           // 43
	xove                                      // 44
	oia_proce_ve                              // 45
	conteudo_string                           // 46
	comentario_linha                          // 47
	causo                                     // 48
	fim_do_causo                              // 49
	conteudo_inteiro                          // 50
	conteudo_hexa                             // 51
	conteudo_octal                            // 52
	conteudo_float                            // 53
	variavel                                  // 54
	)

	var PalavrasReservadas = map[string]TabelaPalavras{
	"trem_di_numeru":      trem_di_numeru,
	"trem_cum_virgula":     trem_cum_virgula,
	"trem_discrita":        trem_discrita,
	"trem_discolhe":        trem_discolhe,
	"trosso":               trosso,
	"uai_se":               uai_se,
	"uai_senao":            uai_senao,
	"roda_esse_trem":       roda_esse_trem,
	"enquanto_tiver_trem":  enquanto_tiver_trem,
	"dependenu":            dependenu,
	"du_casu":              du_casu,
	"ta_bao":               ta_bao,
	"para_o_trem":          para_o_trem,
	"toca_o_trem":          toca_o_trem,
	"bora_cumpade":         bora_cumpade,
	"abre_parentese":       abre_parentese,
	"fecha_parentese":      fecha_parentese,
	"eh":                   eh,
	"num_eh":               num_eh,
	"simbora":              simbora,
	"cabo":                 cabo,
	"abre_chave":           abre_chave,
	"fecha_chave":          fecha_chave,
	"abre_aspas":           abre_aspas,
	"fecha_aspas":          fecha_aspas,
	"virgula":              virgula,
	"uai":                  uai,
	"<":                   	menor_que,
	">":                   	maior_que,
	"<=":                  	menor_igual_que,
	">=":                   maior_igual_que,
	"fica_assim_entao":     fica_assim_entao,
	"neh_nada":             neh_nada,
	"mema_coisa":           mema_coisa,
	"quarque_um":           quarque_um,
	"vam_marca":            vam_marca,
	"tamem":                tamem,
	"um_o_oto":             um_o_oto,
	"+":                    soma,
	"-":                    subtracao,
	"veiz":                 veiz,
	"sob":                  sob,
	"%":               		modulo,
	"/":      				divisao_inteira,
	"xove":                 xove,
	"oia_proce_ve":         oia_proce_ve,
    
	"causo":                causo,
	"fim_do_causo":         fim_do_causo,
	}