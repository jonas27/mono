import 'package:flutter/material.dart';
import 'package:wedding/dayModel.dart';
import 'package:wedding/dayPage.dart';
import 'package:intl/intl.dart';
import 'package:wedding/days/kartoffel.dart';
import 'package:wedding/days/nichts.dart';
import 'package:wedding/days/picknick.dart';
import 'package:wedding/days/movie.dart';
import 'package:wedding/days/wandern.dart';

class OverviewPage extends StatefulWidget {
  const OverviewPage({
    Key key,
  }) : super(key: key);

  @override
  OverviewPageState createState() => OverviewPageState();
}

class OverviewPageState extends State<OverviewPage> {

  @override
  initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    if ( MediaQuery.of(context).size.width <1000){
      return _mobile();
    }

    return Scaffold(
      body: Center(
        child: Row(
          children: [
            Expanded(
              flex: 1,
              child: Container(
                decoration: BoxDecoration(
                  image: DecorationImage(
                    image: AssetImage("assets/images/days/overview.jpg"),
                    fit: BoxFit.cover,
                  ),
                ),
                child: null /* add child content here */,
              ),
            ),
            Expanded(
              flex: 1,
              child: Container(
                color: Color.fromRGBO(6, 28, 48, 1),
                child: ListView(
                  children: [
                    Align(
                      alignment: Alignment.topCenter,
                      child: Padding(
                        padding: const EdgeInsets.fromLTRB(0,60,0,0),
                        child: Text(
                          'Wedding',
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
                        padding: EdgeInsets.all(100.0 ),
                        child: SelectableText(
                          'Hey Katchi, \n\nbald ist es so weit, in 6 Tage heiraten wir und bis dahin gibt es noch 5 Überraschungen für dich. '
                          'Leider spielt das Wetter nicht so mit wie ich mir das erhofft hatte, deswegen werden wir etwas spontaner sein müssen. '
                          'Aber auch das werden wir schaffen 😊\n\n'
                          'Jeden Tag, also um 0:00 Uhr, wird ein neuer Link verfügbar werden. '
                              'Jede Seite hat außerdem noch ein Zitat eines Songs, weißt du welcher? Viel Spaß beim Raten 😘',
//                            textDirection: TextDirection.center,
                          style: TextStyle(
                            fontSize: 28,
                            fontFamily: 'DancingScript',
                            color: Colors.white,
                          ),
                          textAlign: TextAlign.justify,
                        ),
                      ),
                    ),
                    _showLinks(),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  _showLinks() {
    if(DateTime.now().month != 8 && DateTime.now().year == 2020) {
      return Container();
    }
//    int day = 1;
    int day = DateTime.now().day;
    day = day > 5 ? 5 : day;
    List<DayModel> list = <DayModel>[];
    list.add(DayModel(Wandern.WANDERN_TITLE,   Wandern.WANDERN, Wandern.WANDERN_IMG,0,false));
    list.add(DayModel(Movie.MOVIE_TITLE, Movie.MOVIE, Movie.MOVIE_IMG,0,false));
    list.add(DayModel(Picknick.PICKNICK_TITLE, Picknick.PICKNICK, Picknick.PICKNICK_IMG,0,false));
    list.add(DayModel(Kartoffel.KARTOFFEL_TITLE,   Kartoffel.KARTOFFEL, Kartoffel.KARTOFFEL_IMG,0,false));
    list.add(DayModel(Nichts.NICHTS_TITLE,   Nichts.NICHTS, Nichts.NICHTS_IMG,0,false));
    for (int i=0; i<day; i++){
      list[i].active = true;
    }

    return ListView.builder(
        shrinkWrap: true,
        physics: ClampingScrollPhysics(),

//        padding: const EdgeInsets.all(40),
        itemCount: list.length,
        itemBuilder: (BuildContext context, int index) {
            list[index].day = index+1;
          if (list[index].active) {
            return Container(
              height: 50,
              width: 150,
              padding: const EdgeInsets.all(0.0),
              child: MaterialButton(
                materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(200.0),
                    side: BorderSide(color: Colors.transparent)
                ),
                onPressed: () {
                  Navigator.pushNamed(
                    context,
                    DayPage.routeName,
                    arguments: list[index],
//                    MaterialPageRoute(settings: const RouteSettings(name: '/form'), builder: (context) => DayPage(),),
                  );
                },
                child: Text(
                  'Day ' + list[index].day.toString() + ' - ' + list[index].title,
                  style: TextStyle(
                    fontSize: 28,
                    fontFamily: 'DancingScript',
                    color: Colors.white,
                  ),
                  textAlign: TextAlign.justify,
                ),
              ),
            );
          }
          return Container(
            height: 50,
            child: Center(
              child: Text(
                'Day ' + list[index].day.toString(),
                style: TextStyle(
                  fontSize: 28,
                  fontFamily: 'DancingScript',
                  color: Colors.grey,
                ),
                textAlign: TextAlign.justify,
              ),
            ),
          );
        });
  }

  _mobile(){
    return Scaffold(
      body: Center(
        child: Stack(
          children: [Container(
                decoration: BoxDecoration(
                  image: DecorationImage(
                    image: AssetImage("assets/images/days/overview.jpg"),
                    fit: BoxFit.cover,
                  ),
                ),
                child: null /* add child content here */,
              ),
            Padding(
              padding: const EdgeInsets.all(80.0),
              child: SelectableText(
                'Die Seite sieht auf dem Laptop viel besser aus. Also geh auf jonasmanser.com und enjoy!\n\n'
                    'Falls du schon am Laptop bist, mach das Fenster groß.',
//                            textDirection: TextDirection.center,
                style: TextStyle(
                  fontSize: 28,
                  fontFamily: 'DancingScript',
//                  color: Colors.white,
                ),
                textAlign: TextAlign.justify,
                onTap: null,
              ),
            ),

          ],
        ),
      ),
    );
  }
}
