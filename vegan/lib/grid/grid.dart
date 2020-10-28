import 'package:flutter/material.dart';

import '../fact.dart';
import '../globals.dart';

class Grid extends StatelessWidget {
  Grid({Key key, this.width, this.facts}) : super(key: key);

  final List<Fact> facts;
  double width;

  Widget _addCard(double width, Fact fact) {
    return Container(
      width: width * 0.9,
      child: Padding(
        padding: EdgeInsets.only(top: 20.0, left: 0.0, right: 0.0, bottom: 0.0),
        child: Card(
          color: Globals.COLOR_MAIN,
          child: Padding(
            padding: EdgeInsets.only(
                top: 10.0, left: 10.0, right: 10.0, bottom: 10.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.start,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Center(
                  child: Image.asset(fact.img),
                ),
                Center(
                  child: Container(),),
                ExpansionTile(
                  title: Text(fact.title,
                      style: TextStyle(
                          color: Color.fromRGBO(17, 50, 20, 1),
                          fontWeight: FontWeight.bold,
                          fontSize: 20)),
                  children: [
                    for (String text in fact.points)
                      Column(
                        children: [
                          Center(
                            child: Padding(
                              padding: const EdgeInsets.fromLTRB(0, 10, 0, 10),
                              child: Text(
                                text,
                                style: TextStyle(
                                    color: Globals.COLOR_BACKGROUND,
                                    fontSize: 18),
                                textAlign: TextAlign.center,
                              ),
                            ),
                          ),
                          Container(
                            width: width * 0.7,
                            height: 0.8,
                            color: Globals.COLOR_BACKGROUND,
                          )
                        ],
                      ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildGrid({int columns = 3}) {
    List<Widget> rows = [];
    List<Widget> column = [];
    width > 1200 ? width = 1200 : width = width;
    if (width < 600) {
      for (int i = 0; i < facts.length; i++) {
        print(i);
        column.add(_addCard(width, facts[i]));
      }
      return Column(
        children: column,
      );
    } else {
      for (int i = 0; i < columns; i++) {
        for (int j = i; j < facts.length; j = j + 3) {
          column.add(_addCard(width / 3, facts[j]));
        }
        rows.add(Expanded(child: Column(children: column)));
        column = [];
      }
    }
    return Row(
      children: rows,
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisAlignment: MainAxisAlignment.center,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 10, 0, 30),
      child: _buildGrid(),
    );
  }
}
