import 'package:link/link.dart';
import 'package:flutter/material.dart';

import '../fact.dart';
import '../globals.dart';
import '../grid/grid.dart';
import '../youtube.dart';
import 'facts.dart';

class Gesundheit extends StatelessWidget {
  Gesundheit({Key key}) : super(key: key);

  EdgeInsetsGeometry paddingYoutube = EdgeInsets.fromLTRB(0, 30, 0, 30);

  Widget _showRefs(double width) {
    return Column(children: [
      for (Fact fact in Facts.facts)
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0,8,0,8.0),
              child: Text(
                fact.title,
                style: TextStyle(fontSize: 24, color: Globals.COLOR_MAIN),
              ),
            ),
            for (Reference ref in fact.references)
              Padding(
                padding: const EdgeInsets.fromLTRB(30, 0, 0, 0),
                child: Column(crossAxisAlignment: CrossAxisAlignment.start,children: [
                  _customText(ref.author),
                  Row(
                    children: [
                      Container(width: 170,child: _customText(ref.date)),
                      _customText(ref.title),
                    ],
                  ),
                  Link(
                    child: Text(ref.link, style: TextStyle(decoration: TextDecoration.underline, color: Colors.blue)),
                    url: ref.link,
                    onError: null,
                  ),
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Container(width: width,height: 1,color: Globals.COLOR_MAIN,),
                  )
                ]),
              ),
          ],
        ),
    ]);
  }

  Widget _customText(String text){
    return Text(text,style: TextStyle(fontSize: 16, color: Globals.COLOR_MAIN),);
  }

  @override
  Widget build(BuildContext context) {
    double width= MediaQuery.of(context).size.width;
    return ListView(
      children: [
        Center(
          child: Container(
              padding: paddingYoutube,
              child: Youtube(
                'PyejhJeInzQ',
                width: 1200,
                height: MediaQuery.of(context).size.height / 2,
              )),
        ),
        Grid(width: MediaQuery.of(context).size.width, facts: Facts.facts),
        width>900?
        _showRefs(MediaQuery.of(context).size.height / 2.5):Container(),
      ],
    );
  }
}
