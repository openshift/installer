all: install_upi.svg install_upi_vpc.svg

install_upi.svg: install_upi.dia relink-dia.py AWS-Architecture-Icons_SVG/.touch
	dia --export $@ $<
	./relink-dia.py $@

install_upi_vpc.svg: install_upi.dia relink-dia.py AWS-Architecture-Icons_SVG/.touch
	dia -L 'Availability Zones,VPC,Subnets,Route53 (Private)' --export $@ $<
	./relink-dia.py $@

# URI from https://aws.amazon.com/architecture/icons/
AWS-Architecture-Icons_SVG.zip:
	curl https://d1.awsstatic.com/webteam/architecture-icons/AWS-Architecture-Icons_SVG_20190131.f8dacae00681382ab91d59994de429d6e9a680bc.zip >$@

AWS-Architecture-Icons_SVG/.touch: AWS-Architecture-Icons_SVG.zip
	unzip -o $<
	touch $@

clean:
	rm -rf install_upi*.svg AWS-Architecture-Icons_SVG* __MACOSX
