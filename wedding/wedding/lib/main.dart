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

      home: OverviewPage());
  }
  }

