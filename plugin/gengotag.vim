if exists('g:gengotag_loaded')
  finish
endif
let g:gengotag_loaded = 1

let s:basepath = expand('<sfile>:p:h:h')
let s:bin = s:basepath . '/' . 'gengotag'

if !executable(s:bin)
  echom 'Please build `gengotag` first'
  finish
endif

function! s:gen(file, type, omitempty)
  let omit = ""
  if a:omitempty
    let omit = "-omitempty"
  endif
  let cmd = printf('%s -file %s -type %s %s 2>/dev/null', s:bin, a:file, a:type, omit)
  let str = system(cmd)
  let save = @a
  let @a = str
  normal "ap
  let @a = save
endfunction

command! -bang -nargs=1 -complete=file Gengotag call <SID>gen(<q-args>, 'json', "<bang>" == "!")
