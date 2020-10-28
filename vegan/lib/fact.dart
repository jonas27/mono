class Fact {
  final String img;
  final String title;
  final List<String> points;
  final List<Reference> references;

  Fact({this.img='', this.title='', this.points=const [], this.references=const []});
}

class Reference {
  final String author;
  final String title;
  final String link;
  final String date;

  Reference({this.author='', this.link='', this.date='', this.title=''});
}
