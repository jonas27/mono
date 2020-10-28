import 'dart:html' as html;
import 'dart:ui' as ui;

import 'package:flutter/material.dart';

class Vegan extends StatelessWidget {
  Vegan({Key key}) : super(key: key);

  String viewID = "your-view-id";

  @override
  Widget build(BuildContext context) {

    return ListView(
      children: [
        Container(
          height: 20,
        ),
        Container(
          height: 50,
        ),
        Container(constraints: BoxConstraints(maxHeight: 500, maxWidth: 500,),child: Image(image: AssetImage("assets/img/logo/veganismus_colored5.png"))),
      ],
    );
  }
}
