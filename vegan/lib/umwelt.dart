import 'package:flutter/material.dart';

import 'youtube.dart';

class Umwelt extends StatelessWidget {
  Umwelt({Key key}) : super(key: key);

  EdgeInsetsGeometry paddingYoutube = EdgeInsets.fromLTRB(0, 30, 0, 30);


  @override
  Widget build(BuildContext context) {
    return ListView(
      children: [
        Center(
          child: Container(
              padding: paddingYoutube,
              child: Youtube(
                'aETNYyrqNYE',
                width: 1200,
                height: MediaQuery
                    .of(context)
                    .size
                    .height / 2,
              )),
        ),
      ],
    );
  }
}
