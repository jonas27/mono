import 'package:flutter/material.dart';
import 'package:vegan/globals.dart';
import 'gesundheit/gesundheit.dart';
import 'umwelt.dart';
import 'vegan.dart';
import 'ethik.dart';

class Overview extends StatefulWidget {
  Overview({Key key, this.titles, this.icon}) : super(key: key);

  final Icon icon;
  final List<String> titles;

  @override
  _OverviewState createState() => _OverviewState();
}

class _OverviewState extends State<Overview> with SingleTickerProviderStateMixin {
  TabController _tabController;

  // Color colorText = Colors.white70;

  @override
  void initState() {
    super.initState();
    _tabController =
        new TabController(vsync: this, length: widget.titles.length + 1, initialIndex: 3);
  }

  _buildTabs() {
    List<Widget> list = [];
    list.add(Container( constraints: BoxConstraints(maxWidth: 40, maxHeight: 40),child: Image(image: AssetImage("assets/img/logo/veganismus.png"))));
    list.addAll(widget.titles.map((name) => Tab(text: name)).toList());
    return list;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        color: Globals.COLOR_BACKGROUND,
        child: Center(
          child: Container(
            constraints: BoxConstraints(maxWidth: 1200),
            child: DefaultTabController(
                length: widget.titles.length + 1,
                child: Column(
                  children: [
                    TabBar(
                      controller: _tabController,
                      tabs: _buildTabs(),
                      indicatorColor: Globals.COLOR_MAIN,
                      labelColor: Globals.COLOR_MAIN,
                      labelStyle: TextStyle(fontSize: 14,fontWeight: FontWeight.bold),
                      unselectedLabelStyle: TextStyle(fontSize: 14 ),
                      onTap: (index) {
                        // Tab index when user select it, it start from zero
                      },
                    ),
                    Expanded(
                      child: TabBarView(
                          controller: _tabController,
                      children: [
                        Vegan(),
                        Umwelt(),
                        Ethik(),
                        Gesundheit(),
                        // Container(),
                      ]),
                    ),
                  ],
                )),
          ),
        ),
      ),
    );
  }
}
