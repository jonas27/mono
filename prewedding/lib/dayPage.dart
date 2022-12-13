import 'package:flutter/material.dart';
import 'package:wedding/dayModel.dart';

class DayPage extends StatelessWidget {
  static const routeName = '/day';
  final String title;
  final String mainText;
  final String img;
  final int day;
  final bool active;


  const DayPage({
    Key key,
    this.title,
    this.mainText,
    this.img,
    this.day,
    this.active
  }) : super(key: key);


  @override
  Widget build(BuildContext context) {

    return Scaffold(
      body: Center(
        child: Row(
          children: [
            Expanded(
              flex: 1,
              child: Container(
                color: Color.fromRGBO(6, 28, 48, 1),
                child: ListView(
                  children: [
                    Align(
                      alignment: Alignment.topCenter,
                      child: Padding(
                        padding: const EdgeInsets.fromLTRB(0,60,0,60),
                        child: Text(
                          title,
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 46,
                            fontFamily: 'DancingScript',
                            color: Colors.white,
                          ),
                        ),
                      ),
                    ),
                    Align(
                      alignment: Alignment.topCenter,
                      child: Padding(
                        padding: EdgeInsets.fromLTRB(100,40,100,0),
                        child: SelectableText(mainText,
                          style: TextStyle(
                            fontSize: 28,
                            fontFamily: 'DancingScript',
                            color: Colors.white,
                          ),
                          textAlign: TextAlign.justify,
                        ),
                      ),
                    ),
                    _showRings(day),
                  ],
                ),
              ),
            ),
            Expanded(
              flex: 1,
              child: Container(
                decoration: BoxDecoration(
                  image: DecorationImage(
                    image: AssetImage("assets/images/days/" + img),
                    fit: BoxFit.cover,
                  ),
                ),
                child: null /* add child content here */,
              ),
            ),
          ],
        ),

//            child: DayPage(),
      ),
    );
  }

  _showRings(int i){
      if (i==5) {
        return Padding(
          padding: const EdgeInsets.all(60.0),
          child: ConstrainedBox(
            constraints: BoxConstraints(
              maxHeight: 120, maxWidth: 120,
            ),
            child: Opacity(opacity: 0.8,child: Image( image: AssetImage("assets/images/days/ringe.png"),))
          ),
        );
      }
      return Container();
  }
}
