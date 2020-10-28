import 'dart:html' as html;
import 'dart:ui' as ui;

import 'package:flutter/material.dart';

class Youtube extends StatelessWidget {
  Youtube(this.viewID, {Key key, this.width, this.height}) : super(key: key);

  final String viewID;
  final double width;
  final double height;

  @override
  Widget build(BuildContext context) {
    // ignore: undefined_prefixed_name
    ui.platformViewRegistry.registerViewFactory(
        viewID,
        (int id) => html.IFrameElement()
          // ..width = MediaQuery.of(context).size.width.toString()
          ..width = '100'
          ..height = '100'
          // ..height = MediaQuery.of(context).size.height.toString()
          ..src = 'https://www.youtube.com/embed/$viewID'
          ..style.border = 'none');

    return Container(
      // maxHeight: MediaQuery
      //     .of(context)
      //     .size
      //     .height / 2,
      width: width,
      height: height,
      child: HtmlElementView(
        viewType: viewID,
      ),
    );
  }
}
