import 'package:eatat/overviewPage.dart';
import 'package:flutter/material.dart';

void main() {
  runApp(MyApp()
  );
}

class MyApp extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Katja and Jonas 2020',
//      initialRoute: '/',
//      onGenerateRoute: (settings) {
//        if (settings.name == DayPage.routeName) {
//          // Cast the arguments to the correct type: ScreenArguments.
//          final DayModel args = settings.arguments;
//
//          // Then, extract the required data from the arguments and
//          // pass the data to the correct screen.
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

