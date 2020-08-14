import 'package:flutter/cupertino.dart';
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
  List<String> pictureUris=[];

  @override
  initState() {
    for(int i=5; i<19; i++) {
      pictureUris.add("assets/img/wed/" + i.toString() + ".jpg");
    }
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Container(
          color: Color.fromRGBO(6, 28, 48, 1),
          child: Center(
            child: Container(
              constraints: BoxConstraints(
                maxWidth: 960,
              ),
              child: ListView(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 20, 0, 0),
                  ),
//                  Row(
//                    mainAxisAlignment: MainAxisAlignment.end,
//                    children: [
//                      SizedBox(
//                        width: 60,
//                        child: FlatButton(
//                          onPressed: (() {
//                            print("de");
//                          }),
//                          child: Container(
//                            child: Text(
//                              'de',
//                              style: TextStyle(
//                                fontSize: 20,
//                                fontFamily: 'Z003',
//                                color: Colors.grey,
//                              ),
//                            ),
//                          ),
//                        ),
//                      ),
//                      SizedBox(
//                        width: 60,
//                        child: FlatButton(
//                          onPressed: (() {}),
//                          child: Container(
//                            child: Text(
//                              'en',
//                              style: TextStyle(
//                                fontSize: 20,
//                                fontFamily: 'Z003',
//                                color: Colors.grey,
//                              ),
//                            ),
//                          ),
//                        ),
//                      ),
//                    ],
//                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        'Mr ',
                        style: TextStyle(
                          fontSize: 70,
                          fontFamily: 'Z003',
                          color: Colors.white,
                        ),
                      ),
                      Text(
                        '&',
                        style: TextStyle(
                          fontSize: 55,
                          fontFamily: 'MathJax_Caligraphy',
                          color: Colors.white,
                        ),
                      ),
                      Text(
                        ' Mrs',
                        style: TextStyle(
                          fontSize: 70,
                          fontFamily: 'Z003',
                          color: Colors.white,
                        ),
                      ),
                    ],
                  ),
//                  ),
                  Padding(
                    padding: EdgeInsets.fromLTRB(30, 30, 30, 30.0),
                    child: SelectableText(
                      'Wir haben uns entschlossen klammheimlich am 06.08.2020 zu heiraten. '
                      'Das Datum hat für uns eine besondere Bedeutung, da wir uns an diesem Tag vor genau 2 Jahren im Zug nach Copenhagen kennengelernt haben.\n\n'
                      'Jetzt sind wir den nächsten Schritt gegangen, ganz alleine, wollen euch aber auf diesem Weg an unserem Glück teilhaben lassen. '
                      'Deshalb findet ihr hier einige Bilder unseres besonderen Tages.',
                      style: TextStyle(
                        fontSize: 26,
                        fontFamily: 'DancingScript',
                        color: Colors.white,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  Gallery(
                    pictureUris: pictureUris,
                  ),
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Center(
                      child: Column(
                        children: [
                          Opacity(
                            opacity: 0.7,
                            child: Icon(
                              Icons.favorite,
                              color: Colors.pink,
                              size: 40.0,
                            ),
                          ),
                          Text(
                            'By Katja and Jonas Manser, 2020',
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
        ),
      ),
    );
  }
}
