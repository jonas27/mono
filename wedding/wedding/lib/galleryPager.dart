import 'package:flutter/material.dart';

class GalleryPager extends StatefulWidget {
  const GalleryPager({
    Key key,
    this.pictureUris,
    this.initialItem,
  }) : super(key: key);
  final List<String> pictureUris;
  final int initialItem;

  _GalleryPagerState createState() => _GalleryPagerState();
}

class _GalleryPagerState extends State<GalleryPager> {
  PageController _pageController;

  @override
  void initState() {
    super.initState();
    _pageController = PageController(initialPage: widget.initialItem);
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  _createImageList() {
    List<Widget> list = [];
    for (int i = 0; i < widget.pictureUris.length; i++) {
      list.add(
        Stack(
          children: [
            Center(
              child: Image(image: AssetImage(widget.pictureUris[i])),
            ),
            Row(
              children: [
                Expanded(
                  flex: 1,
                  child: GestureDetector(
                    onTap: (() {
                      if (_pageController.hasClients) {
                        _pageController.animateToPage(
                          i - 1,
                          duration: const Duration(milliseconds: 400),
                          curve: Curves.easeInOut,
                        );
                      }
                    }),
                    onVerticalDragEnd: ((DragEndDetails) {
                      Navigator.pop(context);
                    }),
                  ),
                ),
                Expanded(
                  flex: 1,
                  child: GestureDetector(
                    onTap: (() {
                      if (_pageController.hasClients) {
                        _pageController.animateToPage(
                          i + 1,
                          duration: const Duration(milliseconds: 400),
                          curve: Curves.easeInOut,
                        );
                      }
                    }),
                    onVerticalDragEnd: ((DragEndDetails details) {
                      Navigator.pop(context);
                    }),
                  ),
                ),
              ],
            ),
          ],
        ),
      );
    }
    return list;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        color: Color.fromRGBO(6, 28, 48, 1),
        child: PageView(
          controller: _pageController,
          children: _createImageList().toList(),

//        children: [
//          Container(
//            color: Colors.red,
//            child: Center(
//              child: RaisedButton(
//                color: Colors.white,
//                onPressed: () {
//                  if (_pageController.hasClients) {
//                    _pageController.animateToPage(
//                      1,
//                      duration: const Duration(milliseconds: 400),
//                      curve: Curves.easeInOut,
//                    );
//                  }
//                },
//                child: Text('Next'),
//              ),
//            ),
//          ),
//        ],
        ),
      ),
    );
  }
}
