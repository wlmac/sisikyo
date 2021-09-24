pages: events.ics index.html
	mkdir -p public
	cp events.ics public

index.html: index.xhtml
	cp index.xhtml index.html

events.ics:
	go build -o ical gitlab.com/mirukakoro/sisikyo/events/cmd/ical
	./ical > events.ics

