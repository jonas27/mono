import 'package:flutter/material.dart';
import '../http/vocab.dart';

class MyDrawer extends StatelessWidget {
  final bool loggedIn;

  const MyDrawer({Key? key, required this.loggedIn}): super(key: key);

  @override
  Widget build(BuildContext context) {
  return Drawer(
          child: ListView(
            padding: EdgeInsets.zero,
            children: [
              const DrawerHeader(
                decoration: BoxDecoration(
                  color: Colors.blue,
                ),
                child: Text('Burger Menu'),
              ),
              if (loggedIn) ListTile(
                title: const Text('Sign out'),
                onTap: () {
                  getVocabList();
                },
              ),
              if (!loggedIn) ListTile(
                title: const Text('Sign in'),
                onTap: () {
                  // do something
                },
              ),
            ],
          ),
        );
  }
}