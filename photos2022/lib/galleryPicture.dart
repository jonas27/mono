import 'package:flutter/material.dart';
import 'package:photos2022/overviewPage.dart';
import 'package:sprintf/sprintf.dart';

class GalleryPicture extends StatelessWidget {
  const GalleryPicture({
    Key? key,
    required this.pictureWidth,
    required this.pictureHeight,
    required this.num,
  }) : super(key: key);

  final double pictureWidth;
  final double pictureHeight;
  final int num;

  @override
  Widget build(BuildContext context) {
    String uri = sprintf(URI_MINI, [num]);
    precacheImage(AssetImage(uri), context);
    return AspectRatio(
      aspectRatio: 1 / 1,
      child: Container(
        decoration: BoxDecoration(
          shape: BoxShape.rectangle,
          image: DecorationImage(
            image: AssetImage(uri),
            fit: BoxFit.cover,
          ),
        ),
        child: null /* add child content here */,
      ),
    );
  }
}
