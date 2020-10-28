import '../fact.dart';

class Facts {
  static final List<Fact> facts = [
    Fact(img: 'assets/img/gesundheit/herz.jpg', title: 'Krebs', points: [
      'Krebs ist weltweit die zweithäufigste Ursache für Tod. Etwa jede sechste stirbt an Krebs. Das sind 9.6 Mio Menschen weltweit',
      'Lifestyleveränderungen, wie whole food plant-based diet, können Prostatekrebswachstum um das 8-fache hemmen. Für andere Krebsarten ist es wahrscheinlich ähnlich.',
      'Rund ein Drittel der Krebstote weltweit sind auf die fünf wichtigsten Verhaltens- und Ernährungsrisiken zurückzuführen.',
      'Ursachen in Deutschland: Rauchen 19.3%, Ungesunde Ernährung: 7.6%, Übergewicht: 6.9%, Bewegungsmangel: 6.1%, Alkohol: 2.2%.',
      'In Deutschland sterben jährlich ca. 230.000 an Krebs.',
      'Die wirtschaftlichen Auswirkungen von Krebs sind erheblich und nehmen zu. Die jährlichen global Gesamtkosten sind im Jahr 2010 ca. 1,16 Billionen USD',
      'In Deutschland sind 37% der 440.000 aller neuen Krebsfälle vermeidbar. Das sind ca. 163.000 vermeidbare Krebsfälle.',
    ], references: [
      Reference(
          author: 'Bundesministerium für Gesundheit',
          link:
              'https://www.bundesgesundheitsministerium.de/themen/praevention/gesundheitsgefahren/krebs.html',
          date: '23. Oktober 2020',
          title: 'Krebs'),
      Reference(
          author: 'WHO',
          link: 'https://www.who.int/news-room/fact-sheets/detail/cancer',
          date: '12 September 2018',
          title: 'Cancer'),
      Reference(
          author:
              'Frattaroli J, Weidner G, Dnistrian AM, Kemp C, Daubenmier JJ, Marlin RO, Crutchfield L, Yglecias L, Carroll PR, Ornish D',
          link: 'https://pubmed.ncbi.nlm.nih.gov/18602144/',
          date: '12 September 2018',
          title:
              'Clinical events in prostate cancer lifestyle trial: results from two years of follow-up.'),
      Reference(
          author: 'Deutschen Krebsforschungszentrum',
          link:
              'https://www.dkfz.de/de/presse/pressemitteilungen/2018/dkfz-pm-18-48-Vermeidbare-Risikofaktoren-verursachen-37-Prozent-aller-Krebsfaelle.php',
          date: '03.09.2018',
          title:
              'Erstmals für Deutschland ermittelt: Vermeidbare Risikofaktoren verursachen 37 Prozent aller Krebsfälle'),
    ]),
  ];
}
