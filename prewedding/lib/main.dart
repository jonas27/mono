import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:wedding/dayModel.dart';
import 'package:wedding/dayPage.dart';
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
//      routes: {
//        DayPage.routeName: (context) => DayPage(),
//      },
//      onGenerateRoute: (settings) {
//        final arguments = settings.arguments;
//        switch (settings.name) {
//          case '/day':
//            if (arguments is String) {
//              // the details page for one specific user
//              return DayModel(arguments);
//            }
//            else {
//              // a route showing the list of all users
//              return UserList();
//            }
//            break;
//          default:
//            return null;
//        }
        onGenerateRoute: (settings) {
      // If you push the PassArguments route
      if (settings.name == DayPage.routeName) {
        // Cast the arguments to the correct type: ScreenArguments.
        final DayModel args = settings.arguments;

        // Then, extract the required data from the arguments and
        // pass the data to the correct screen.
        return MaterialPageRoute(
            builder: (context) {
          return DayPage(
            title: args.title,
            mainText: args.mainText,
              img: args.img,
              day: args.day,
              active: args.active

          );
        },
        );}
  },
        home: OverviewPage(),);
}
}

