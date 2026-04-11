package prayer

import "time"

// DuaEntry represents a single du'a with its text and source
type DuaEntry struct {
	Arabic          string
	Transliteration string
	Translation     string
	Reference       string
}

// Duas contains a collection of du'as from Hisn al-Muslim and Quran
var Duas = []DuaEntry{
	{
		Arabic:          "رَبَّنَا آتِنَا فِي الدُّنْيَا حَسَنَةً وَفِي الْآخِرَةِ حَسَنَةً وَقِنَا عَذَابَ النَّارِ",
		Transliteration: "Rabbana atina fid-dunya hasanatan wa fil-akhirati hasanatan waqina 'adhab an-nar",
		Translation:     "Our Lord, give us good in this world and good in the Hereafter, and protect us from the torment of the Fire",
		Reference:       "Quran 2:201",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَسْأَلُكَ الْعَفْوَ وَالْعَافِيَةَ فِي الدُّنْيَا وَالْآخِرَةِ",
		Transliteration: "Allahumma inni as'aluka al-'afwa wal-'afiyah fid-dunya wal-akhirah",
		Translation:     "O Allah, I ask You for pardon and well-being in this life and the next",
		Reference:       "Ibn Majah",
	},
	{
		Arabic:          "اللَّهُمَّ أَعِنِّي عَلَى ذِكْرِكَ وَشُكْرِكَ وَحُسْنِ عِبَادَتِكَ",
		Transliteration: "Allahumma a'inni 'ala dhikrika wa shukrika wa husni 'ibadatik",
		Translation:     "O Allah, help me to remember You, thank You, and worship You well",
		Reference:       "Abu Dawud",
	},
	{
		Arabic:          "رَبِّ اشْرَحْ لِي صَدْرِي وَيَسِّرْ لِي أَمْرِي",
		Transliteration: "Rabbi ishrah li sadri wa yassir li amri",
		Translation:     "My Lord, expand for me my chest and ease for me my task",
		Reference:       "Quran 20:25-26",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَعُوذُ بِكَ مِنَ الْهَمِّ وَالْحَزَنِ",
		Transliteration: "Allahumma inni a'udhu bika minal-hammi wal-hazan",
		Translation:     "O Allah, I seek refuge in You from anxiety and sorrow",
		Reference:       "Bukhari",
	},
	{
		Arabic:          "رَبَّنَا لَا تُزِغْ قُلُوبَنَا بَعْدَ إِذْ هَدَيْتَنَا وَهَبْ لَنَا مِن لَّدُنكَ رَحْمَةً",
		Transliteration: "Rabbana la tuzigh qulubana ba'da idh hadaytana wa hab lana min ladunka rahmah",
		Translation:     "Our Lord, let not our hearts deviate after You have guided us, and grant us mercy from Yourself",
		Reference:       "Quran 3:8",
	},
	{
		Arabic:          "حَسْبُنَا اللَّهُ وَنِعْمَ الْوَكِيلُ",
		Transliteration: "Hasbunallahu wa ni'mal-wakil",
		Translation:     "Sufficient for us is Allah, and He is the best Disposer of affairs",
		Reference:       "Quran 3:173",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَسْأَلُكَ عِلْمًا نَافِعًا وَرِزْقًا طَيِّبًا وَعَمَلًا مُتَقَبَّلًا",
		Transliteration: "Allahumma inni as'aluka 'ilman nafi'an wa rizqan tayyiban wa 'amalan mutaqabbalan",
		Translation:     "O Allah, I ask You for beneficial knowledge, good provision, and accepted deeds",
		Reference:       "Ibn Majah",
	},
	{
		Arabic:          "رَبِّ زِدْنِي عِلْمًا",
		Transliteration: "Rabbi zidni 'ilma",
		Translation:     "My Lord, increase me in knowledge",
		Reference:       "Quran 20:114",
	},
	{
		Arabic:          "اللَّهُمَّ اهْدِنِي وَسَدِّدْنِي",
		Transliteration: "Allahumma-hdini wa saddidni",
		Translation:     "O Allah, guide me and keep me on the right path",
		Reference:       "Muslim",
	},
	{
		Arabic:          "رَبَّنَا اغْفِرْ لَنَا ذُنُوبَنَا وَإِسْرَافَنَا فِي أَمْرِنَا",
		Transliteration: "Rabbana-ghfir lana dhunubana wa israfana fi amrina",
		Translation:     "Our Lord, forgive us our sins and our excesses in our affairs",
		Reference:       "Quran 3:147",
	},
	{
		Arabic:          "اللَّهُمَّ إِنَّكَ عَفُوٌّ تُحِبُّ الْعَفْوَ فَاعْفُ عَنِّي",
		Transliteration: "Allahumma innaka 'afuwwun tuhibbul-'afwa fa'fu 'anni",
		Translation:     "O Allah, You are Most Forgiving, and You love forgiveness, so forgive me",
		Reference:       "Tirmidhi",
	},
	{
		Arabic:          "رَبِّ أَوْزِعْنِي أَنْ أَشْكُرَ نِعْمَتَكَ الَّتِي أَنْعَمْتَ عَلَيَّ",
		Transliteration: "Rabbi awzi'ni an ashkura ni'matakal-lati an'amta 'alayya",
		Translation:     "My Lord, enable me to be grateful for Your favor which You have bestowed upon me",
		Reference:       "Quran 27:19",
	},
	{
		Arabic:          "اللَّهُمَّ بَارِكْ لِي فِيمَا رَزَقْتَنِي وَقِنِي عَذَابَ النَّارِ",
		Transliteration: "Allahumma barik li fima razaqtani waqini 'adhab an-nar",
		Translation:     "O Allah, bless me in what You have provided me and protect me from the punishment of the Fire",
		Reference:       "Hisn al-Muslim",
	},
	{
		Arabic:          "اللَّهُمَّ اجْعَلْنِي مِنَ التَّوَّابِينَ وَاجْعَلْنِي مِنَ الْمُتَطَهِّرِينَ",
		Transliteration: "Allahumma-j'alni minat-tawwabina waj'alni minal-mutatahhirin",
		Translation:     "O Allah, make me among those who repent and make me among those who purify themselves",
		Reference:       "Tirmidhi",
	},
	{
		Arabic:          "لَا إِلَٰهَ إِلَّا أَنتَ سُبْحَانَكَ إِنِّي كُنتُ مِنَ الظَّالِمِينَ",
		Transliteration: "La ilaha illa anta subhanaka inni kuntu minaz-zalimin",
		Translation:     "There is no deity except You; exalted are You. Indeed, I have been of the wrongdoers",
		Reference:       "Quran 21:87",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَسْأَلُكَ الْهُدَى وَالتُّقَى وَالْعَفَافَ وَالْغِنَى",
		Transliteration: "Allahumma inni as'alukal-huda wat-tuqa wal-'afafa wal-ghina",
		Translation:     "O Allah, I ask You for guidance, piety, chastity, and self-sufficiency",
		Reference:       "Muslim",
	},
	{
		Arabic:          "رَبَّنَا تَقَبَّلْ مِنَّا إِنَّكَ أَنتَ السَّمِيعُ الْعَلِيمُ",
		Transliteration: "Rabbana taqabbal minna innaka antas-Sami'ul-'Alim",
		Translation:     "Our Lord, accept from us. Indeed You are the All-Hearing, the All-Knowing",
		Reference:       "Quran 2:127",
	},
	{
		Arabic:          "اللَّهُمَّ اكْفِنِي بِحَلاَلِكَ عَنْ حَرَامِكَ وَأَغْنِنِي بِفَضْلِكَ عَمَّنْ سِوَاكَ",
		Transliteration: "Allahumma-kfini bihalalika 'an haramika wa aghnini bifadlika 'amman siwak",
		Translation:     "O Allah, suffice me with what is lawful against what is unlawful, and enrich me by Your bounty from all besides You",
		Reference:       "Tirmidhi",
	},
	{
		Arabic:          "رَبَّنَا هَبْ لَنَا مِنْ أَزْوَاجِنَا وَذُرِّيَّاتِنَا قُرَّةَ أَعْيُنٍ",
		Transliteration: "Rabbana hab lana min azwajina wa dhurriyyatina qurrata a'yun",
		Translation:     "Our Lord, grant us from our spouses and offspring comfort to our eyes",
		Reference:       "Quran 25:74",
	},
	{
		Arabic:          "يَا مُقَلِّبَ الْقُلُوبِ ثَبِّتْ قَلْبِي عَلَى دِينِكَ",
		Transliteration: "Ya muqallibal-qulubi thabbit qalbi 'ala dinik",
		Translation:     "O Turner of the hearts, make my heart firm upon Your religion",
		Reference:       "Tirmidhi",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَعُوذُ بِكَ مِنْ عِلْمٍ لَا يَنْفَعُ",
		Transliteration: "Allahumma inni a'udhu bika min 'ilmin la yanfa'",
		Translation:     "O Allah, I seek refuge in You from knowledge that does not benefit",
		Reference:       "Muslim",
	},
	{
		Arabic:          "اللَّهُمَّ أَصْلِحْ لِي دِينِي الَّذِي هُوَ عِصْمَةُ أَمْرِي",
		Transliteration: "Allahumma aslih li dini alladhi huwa 'ismatu amri",
		Translation:     "O Allah, rectify my religion which is the safeguard of my affairs",
		Reference:       "Muslim",
	},
	{
		Arabic:          "رَبِّ اجْعَلْنِي مُقِيمَ الصَّلَاةِ وَمِن ذُرِّيَّتِي",
		Transliteration: "Rabbij-'alni muqimas-salati wa min dhurriyyati",
		Translation:     "My Lord, make me an establisher of prayer, and from my descendants",
		Reference:       "Quran 14:40",
	},
	{
		Arabic:          "اللَّهُمَّ آتِ نَفْسِي تَقْوَاهَا وَزَكِّهَا أَنْتَ خَيْرُ مَنْ زَكَّاهَا",
		Transliteration: "Allahumma ati nafsi taqwaha wa zakkiha anta khayru man zakkaha",
		Translation:     "O Allah, grant my soul its piety and purify it, for You are the best to purify it",
		Reference:       "Muslim",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَسْأَلُكَ الثَّبَاتَ فِي الْأَمْرِ وَالْعَزِيمَةَ عَلَى الرُّشْدِ",
		Transliteration: "Allahumma inni as'alukat-thabata fil-amri wal-'azimata 'alar-rushd",
		Translation:     "O Allah, I ask You for steadfastness in all my affairs and determination upon guidance",
		Reference:       "Ahmad",
	},
	{
		Arabic:          "رَبَّنَا ظَلَمْنَا أَنفُسَنَا وَإِن لَّمْ تَغْفِرْ لَنَا وَتَرْحَمْنَا لَنَكُونَنَّ مِنَ الْخَاسِرِينَ",
		Transliteration: "Rabbana zalamna anfusana wa in lam taghfir lana wa tarhamna lanakoonanna minal-khasirin",
		Translation:     "Our Lord, we have wronged ourselves, and if You do not forgive us and have mercy upon us, we will surely be among the losers",
		Reference:       "Quran 7:23",
	},
	{
		Arabic:          "اللَّهُمَّ إِنِّي أَسْأَلُكَ خَيْرَ هَذَا الْيَوْمِ وَخَيْرَ مَا بَعْدَهُ",
		Transliteration: "Allahumma inni as'aluka khayra hadhal-yawmi wa khayra ma ba'dah",
		Translation:     "O Allah, I ask You for the good of this day and the good of what comes after it",
		Reference:       "Muslim",
	},
	{
		Arabic:          "سُبْحَانَ اللَّهِ وَبِحَمْدِهِ سُبْحَانَ اللَّهِ الْعَظِيمِ",
		Transliteration: "SubhanAllahi wa bihamdihi, SubhanAllahil-'Azim",
		Translation:     "Glory be to Allah and His is the praise, Glory be to Allah the Almighty",
		Reference:       "Bukhari & Muslim",
	},
}

// GetDailyDua returns a du'a that rotates daily based on the date
func GetDailyDua(date time.Time) *DuaEntry {
	if len(Duas) == 0 {
		return nil
	}
	index := date.YearDay() % len(Duas)
	return &Duas[index]
}
