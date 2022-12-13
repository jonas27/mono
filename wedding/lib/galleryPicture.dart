import 'package:flutter/material.dart';


class GalleryPicture extends StatelessWidget {
  const GalleryPicture({
    Key key,
    this.pictureWidth,
    this.pictureHeight,
    this.URI,
  }) : super(key: key);

  final double pictureWidth;
  final double pictureHeight;
  final String URI;

  @override
  Widget build(BuildContext context) {
    return Container(
      child: AspectRatio(
        aspectRatio: 1 / 1,
        child: Container(
          decoration: BoxDecoration(
            shape: BoxShape.rectangle,
            image: DecorationImage(
              image: AssetImage(URI),
              fit: BoxFit.cover,
            ),
          ),
          child: null /* add child content here */,
        ),
      ),
    );
  }
}