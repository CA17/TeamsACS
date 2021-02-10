FROM alpine
COPY release/teamsacs /teamsacs
RUN chmod +x /teamsacs
EXPOSE 1979 1980 1981 1935 1936 1812/udp 1813/udp 1914/udp 1924/udp 1914/udp
ENTRYPOINT ["/teamsacs"]
