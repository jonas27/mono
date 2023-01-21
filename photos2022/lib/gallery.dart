import 'package:flutter/material.dart';
import 'package:photos2022/galleryPager.dart';
import 'package:photos2022/galleryPicture.dart';
import 'package:photos2022/overviewPage.dart';
import 'package:sprintf/sprintf.dart';

class Gallery extends StatefulWidget {
  const Gallery({
    Key? key,
    required this.pictureUris,
    this.picturePerRow = 2,
    this.pictureWidth = 140,
    this.pictureHeight = 140,
  }) : super(key: key);

  final List<int> pictureUris;
  final int picturePerRow;
  final double pictureWidth;
  final double pictureHeight;

  @override
  GalleryState createState() => GalleryState();
}

class GalleryState extends State<Gallery> {
  @override
  initState() {
    super.initState();
  }

  _addImages() {
    List<Widget> list = [];
    for (int i in widget.pictureUris) {
      list.add(GestureDetector(
          onTap: (() {
            _showFullScreen(i);
          }),
          child: GalleryPicture(
            pictureWidth: widget.pictureWidth,
            pictureHeight: widget.pictureHeight,
            num: i,
          )));
    }
    return list;
  }

  _showFullScreen(int i) {
    String uri = sprintf(URI_NORMAL, [i]);
    precacheImage(AssetImage(uri), context);
    Navigator.push(
      context,
      MaterialPageRoute(
          builder: (context) => GalleryPager(
                pictureUris: [uri],
                // pictureUris: widget.pictureUris,
                initialItem: i,
              )),
//      arguments: list[index],
//                    MaterialPageRoute(settings: const RouteSettings(name: '/form'), builder: (context) => DayPage(),),
    );
  }

  @override
  Widget build(BuildContext context) {
    return GridView.count(
        shrinkWrap: true,
        primary: false,
        padding: const EdgeInsets.all(20),
        crossAxisSpacing: 10,
        mainAxisSpacing: 10,
        crossAxisCount: 2,
        children: _addImages().toList());
  }
}
