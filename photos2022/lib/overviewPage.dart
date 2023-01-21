import 'dart:io';
import 'package:flutter/material.dart';
import 'package:photos2022/gallery.dart';

class OverviewPage extends StatefulWidget {
  const OverviewPage({
    Key? key,
  }) : super(key: key);

  @override
  OverviewPageState createState() => OverviewPageState();
}

class OverviewPageState extends State<OverviewPage> {
  final double pictureWidth = 140;
  List<String> pictureUris = [];

  @override
  initState() {
//     Directory dir = Directory('assets/img/wed/');
// // execute an action on each entry
//     dir.list(recursive: false).forEach((f) {
//       print(f.path);
//     });

    // Directory dir = Directory('assets/img/wed/');
    // dir.list(recursive: false, followLinks: false).forEach((element) {
    //   print(element.path);
    // });

    for (int i = 0; i < 68; i++) {
      pictureUris.add("assets/img/photos/${i.toString()}.jpg");
    }
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Container(
          color: const Color.fromRGBO(6, 28, 48, 1),
          child: Center(
            child: Container(
              constraints: const BoxConstraints(
                maxWidth: 960,
              ),
              child: ListView(
                children: [
                  const Padding(
                    padding: EdgeInsets.fromLTRB(0, 20, 0, 0),
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: const [
                      Text(
                        'Photos ',
                        style: TextStyle(
                          fontSize: 70,
                          fontFamily: 'Z003',
                          color: Colors.white,
                        ),
                      ),
                      Text(
                        '2022',
                        style: TextStyle(
                          fontSize: 70,
                          fontFamily: 'Z003',
                          color: Colors.white,
                        ),
                      ),
                    ],
                  ),
//                  ),
                  const Padding(
                    padding: EdgeInsets.fromLTRB(30, 30, 30, 30.0),
                    child: SelectableText(
                      '2022, ein Jahr das irgendwie super schnell vorbei war und in dem es keine wirklich prägenden Momente gab. '
                      'Und gerade deswegen muss man manchmal einfach die kleinen Dinge im Leben lieben... wie dich 😘. \n\n'
                      'Die Zeit mit dir in Copenhagen, Berlin, in Kroatien mit dem Tri-Team oder mit unserer Familie, haben dieses Jahr trotzdem wieder einzigartig gemacht. \n\n'
                      'Natürlich müssen wir die Bilder und den Text noch anpassen und ich muss noch an der Performance arbeiten.',
                      style: TextStyle(
                        fontSize: 26,
                        fontFamily: 'DancingScript',
                        color: Colors.white,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  Gallery(
                    pictureUris: pictureUris,
                  ),
                  Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Center(
                      child: Column(
                        children: const [
                          Opacity(
                            opacity: 0.7,
                            child: Icon(
                              Icons.favorite,
                              color: Colors.pink,
                              size: 40.0,
                            ),
                          ),
                          Text(
                            'By Katja and Jonas Manser, 2022',
                            style: TextStyle(
                              fontSize: 14,
                              fontFamily: 'DancingScript',
                              color: Colors.white,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
