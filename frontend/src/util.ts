export const fancyTimeFormat = (duration: number): string => {   
  // Hours, minutes and seconds
  var hrs = ~~(duration / 3600);
  var mins = ~~((duration % 3600) / 60);
  var secs = ~~duration % 60;

  // Output like "1:01" or "4:03:59" or "123:03:59"
  var ret = "";

  if (hrs > 0) {
      ret += "" + hrs + ":" + (mins < 10 ? "0" : "");
  }

  ret += "" + mins + ":" + (secs < 10 ? "0" : "");
  ret += "" + secs;
  return ret;
}

export const hoursMinutesString = (t: Date): string => {
  let h = t.getHours()
  let ap = 'AM'
  if (h > 12) {
    h -= 12
    ap = 'PM'
  }
  let m = `${t.getMinutes()}`
  if (m.length !== 2) {
    m = '0' + m
  }
  return `${h}:${m} ${ap}`
}