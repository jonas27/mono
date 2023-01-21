import 'package:flutter/material.dart';
import 'package:flutter_countdown_timer/index.dart';

class CountdownTimerPage extends StatefulWidget {
  @override
  _CountdownTimerPageState createState() => _CountdownTimerPageState();
}

class _CountdownTimerPageState extends State<CountdownTimerPage> {
  late CountdownTimerController controller;
  int endTime =
      DateTime(2022, 12, 24).difference(DateTime.now()).inMilliseconds +
          DateTime.now().millisecondsSinceEpoch;

  @override
  void initState() {
    super.initState();
    controller = CountdownTimerController(endTime: endTime);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
        children: <Widget>[
          Center(
            child: CountdownTimer(
              textStyle: TextStyle(
                fontSize: 30,
                color: Colors.red,
              ),
              endTime: endTime,
            ),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }
}
