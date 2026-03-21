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
	abre_aspas                                // 23
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