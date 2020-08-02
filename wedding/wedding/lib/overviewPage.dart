import 'package:flutter/material.dart';
import 'package:wedding/gallery.dart';

class OverviewPage extends StatefulWidget {
  const OverviewPage({
    Key key,
  }) : super(key: key);

  @override
  OverviewPageState createState() => OverviewPageState();
}

class OverviewPageState extends State<OverviewPage> {
  final double pictureWidth = 140;
  List<String> pictureUris = [
  'assets/img/overview.jpg',
  'assets/img/test.png',
  'assets/img/hoch.jpg',
//  'assets/img/overview.jpg',
  'assets/img/overview.jpg',
//  'assets/img/overview.jpg',
//  'assets/img/overview.jpg',
//  'assets/img/overview.jpg',
  ];

  @override
  initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        color: Color.fromRGBO(6, 28, 48, 1),
        child: Center(
          child: Container(
            constraints: BoxConstraints(
              maxWidth: 960,
            ),
            child: Row(
              children: [
                Expanded(
                  child: Container(
//                    color: Color.fromRGBO(6, 28, 48, 1),
                    child: ListView(
                      children: [
                        Align(
                          alignment: Alignment.topCenter,
                          child: Padding(
                            padding: const EdgeInsets.fromLTRB(0, 60, 0, 0),
                            child: RichText(
                              text: TextSpan(
                                children: <TextSpan>[
                                  TextSpan(
                                    text: 'Mr ',
                                    style: TextStyle(
                                      fontSize: 70,
                                      fontFamily: 'Z003',
                                      color: Colors.white,
                                    ),
                                  ),
                                  TextSpan(
                                    text: '&',
                                    style: TextStyle(
                                      fontSize: 55,
                                      fontFamily: 'MathJax_Caligraphy',
                                      color: Colors.white,
                                    ),
                                  ),
                                  TextSpan(
                                    text: ' Mrs',
                                    style: TextStyle(
                                      fontSize: 70,
                                      fontFamily: 'Z003',
                                      color: Colors.white,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                        Padding(
                          padding: EdgeInsets.fromLTRB(30, 30, 30, 30.0),
                          child: SelectableText(
                            'Wir haben uns entschlossen klammheimlich am 06.08.2020 zu heiraten. '
                            'Das Datum hat für uns eine besondere Bedeutung, da wir uns an diesem Tag vor genau 2 Jahren im Zug nach Copenhagen kennengelernt haben.\n\n'
                            'Jetzt sind wir den nächsten Schritt gegangen, ganz alleine, wollen euch aber auf diesem Weg an unserem Glück teilhaben lassen. '
                            'Daher dachten wir, wir teilen mit euch die Bilder unseres besonderen Tages.',
                            style: TextStyle(
                              fontSize: 26,
                              fontFamily: 'DancingScript',
                              color: Colors.white,
                            ),
                            textAlign: TextAlign.center,
                          ),
                        ),
                        Gallery(pictureUris: pictureUris,),
                        Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Center(
                            child: Column(
                              children: [
                              Opacity(opacity: 0.7,child:
                                Icon(
                                  Icons.favorite,
                                  color: Colors.pink,
                                  size: 40.0,
                                ),),
                                Text('By Katja and Jonas Manser, 2020',
                                  style: TextStyle(
                                  fontSize: 14,
                                  fontFamily: 'DancingScript',
                                  color: Colors.white,
                                ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
