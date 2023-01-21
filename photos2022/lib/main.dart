import 'package:flutter/material.dart';
import 'package:photos2022/overviewPage.dart';
import 'package:photos2022/counterPage.dart';
//import 'package:flutter_localizations/flutter_localizations.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    final now = DateTime.now();
    final xmas = DateTime(2022, 12, 24);
    if (now.compareTo(xmas) < 0) {
      return MaterialApp(
        title: 'Flutter Demo',
        theme: ThemeData(
          primarySwatch: Colors.blue,
          visualDensity: VisualDensity.adaptivePlatformDensity,
        ),
        home: MyHomePage(),
      );
    } else {
      return const MaterialApp(
          title: 'Katja and Jonas 2022',
//        initialRoute: '/',
          home: OverviewPage());
    }
  }
}

class MyHomePage extends StatefulWidget {
  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  @override
  Widget build(BuildContext context) {
    return CountdownTimerPage();
  }
}
