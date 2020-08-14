import 'package:flutter/material.dart';

class OverviewPage extends StatefulWidget {
  const OverviewPage({
    Key key,
  }) : super(key: key);

  @override
  OverviewPageState createState() => OverviewPageState();
}

class OverviewPageState extends State<OverviewPage> {

  final double fontsize = 34;

  @override
  initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: Center(
      child: Container(
        constraints: BoxConstraints(maxWidth: 960),
//        color: Color.fromRGBO(6, 28, 48, 1),
        child: ListView(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 60, 0, 40),
              child: Text(
                'Speisekarte',
                style: TextStyle(
//                    fontWeight: FontWeight.bold,
                  fontSize: 46,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Padding(
              padding: EdgeInsets.fromLTRB(20,20,20,40.0),
              child: SelectableText(
                'Das Menu für Familie Burster, wenn sie zu Familie Manser zum Essen kommen will.\n\n'
                    'Stay tuned...',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: 28,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Padding(
              padding: EdgeInsets.fromLTRB(20,20,20,20.0),
              child: SelectableText(
                'Vorspeise',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: fontsize,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Padding(
              padding: EdgeInsets.fromLTRB(20,20,20,20.0),
              child: SelectableText(
                '***',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: fontsize,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Padding(
              padding: EdgeInsets.fromLTRB(20,20,20,20.0),
              child: SelectableText(
                'Hauptgericht',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: fontsize,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Padding(
              padding: EdgeInsets.fromLTRB(20,20,20,20.0),
              child: SelectableText(
                '***',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: fontsize,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
            Padding(
              padding: EdgeInsets.fromLTRB(20,20,20,20.0),
              child: SelectableText(
                'Desert',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: fontsize,
                  fontFamily: 'DancingScript',
                  color: Colors.black,
                ),
                textAlign: TextAlign.center,
              ),
            ),
          ],
        ),
      ),
    ));
  }
}
