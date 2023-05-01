import 'package:flutter/material.dart';

class VocabList extends StatelessWidget {
  const VocabList({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    List<int> ints = [];
    for (var  i=0;i<100;i++) {
      ints.add(i);
    }
    return ListView(
      padding: const EdgeInsets.all(20),
      children: <Widget>[
        for (var item in ints)
        Container(
          height: 50,
          color: Colors.amber[600],
          child: const Center(child: Text('Entry A')),
        ),
        Container(
          height: 50,
          color: Colors.amber[500],
          child: const Center(child: Text('Entry B')),
        ),
        Container(
          height: 50,
          color: Colors.amber[100],
          child: const Center(child: Text('Entry C')),
        ),
      ],
    );
  }
}
