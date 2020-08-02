import 'package:flutter/material.dart';
import 'package:wedding/overviewPage.dart';
import 'package:flutter/foundation.dart' show kIsWeb;

void main() {
  runApp(MyApp()
  );
}

class MyApp extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Katja and Jonas 2020',
      initialRoute: '/',
      onGenerateRoute: (settings) {
//        if (settings.name == DayPage.routeName) {
//          final DayModel args = settings.arguments;
//          return MaterialPageRoute(
//            builder: (context) {
//              return DayPage(
//                  title: args.title,
//                  mainText: args.mainText,
//                  img: args.img,
//                  day: args.day,
//                  active: args.active
//
//              );
//            },
//          );}
//      },
      home: OverviewPage(),);
  }
}

