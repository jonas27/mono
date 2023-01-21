import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class GalleryPager extends StatefulWidget {
  const GalleryPager({
    Key? key,
    required this.pictureUris,
    required this.initialItem,
  }) : super(key: key);
  final List<String> pictureUris;
  final int initialItem;

  @override
  GalleryPagerState createState() => GalleryPagerState();
}

class GalleryPagerState extends State<GalleryPager> {
  late PageController _pageController;
  FocusNode _focusNode = FocusNode();

  @override
  void initState() {
    super.initState();
    _pageController = PageController(initialPage: widget.initialItem);
  }

  @override
  void dispose() {
    // _pageController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  _createImageList() {
    _focusNode = FocusNode();
    List<Widget> list = [];
    for (int i = 0; i < widget.pictureUris.length; i++) {
      list.add(
        RawKeyboardListener(
          autofocus: true,
          focusNode: _focusNode,
          onKey: ((event) {
            if (event.logicalKey == LogicalKeyboardKey.arrowRight ||
                event.physicalKey == PhysicalKeyboardKey.arrowRight) {
              if (_pageController.hasClients) {
                _pageController.animateToPage(
                  i + 1,
                  duration: const Duration(milliseconds: 400),
                  curve: Curves.easeInOut,
//                _focusNode.attach(context);
                );
              }
            } else if (event.logicalKey == LogicalKeyboardKey.arrowLeft ||
                event.physicalKey == PhysicalKeyboardKey.arrowLeft) {
              if (_pageController.hasClients) {
                _pageController.animateToPage(
                  i - 1,
                  duration: const Duration(milliseconds: 400),
                  curve: Curves.easeInOut,
                );
              }
            }
          }),
          child: Stack(
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
        ),
      );
    }
    return list;
  }

  _handleKeyEvent(RawKeyEvent event, int i) {
    print('pressed');
    if (event.logicalKey == LogicalKeyboardKey.arrowRight ||
        event.physicalKey == PhysicalKeyboardKey.arrowRight) {
      if (_pageController.hasClients) {
        _pageController.animateToPage(
          i + 1,
          duration: const Duration(milliseconds: 400),
          curve: Curves.easeInOut,
        );
      }
    } else if (event.logicalKey == LogicalKeyboardKey.arrowLeft ||
        event.physicalKey == PhysicalKeyboardKey.arrowLeft) {
      if (_pageController.hasClients) {
        _pageController.animateToPage(
          i - 1,
          duration: const Duration(milliseconds: 400),
          curve: Curves.easeInOut,
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        color: Color.fromRGBO(6, 28, 48, 1),

//          RawKeyboardListener(
//            autofocus: false,
//            focusNode: FocusNode(),
//            onKey: ((event) {
//              print(event.data.logicalKey.keyId);
//              print('sdf');
//              if (event.runtimeType == PhysicalKeyboardKey.arrowLeft) {
//                bool shiftPressed = event.isShiftPressed; // true: if shift key is pressed
//              }
//            }) ,

        child: PageView(
          controller: _pageController,
          children: _createImageList().toList(),
        ),
      ),
    );
  }
}
